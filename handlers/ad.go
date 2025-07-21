package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vtask/database"
	"vtask/internal/ad"
)

func CreateAd(w http.ResponseWriter, r *http.Request) {
   userIDVal := r.Context().Value("user_id")
    if userIDVal == nil {
        http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
        return
    }
    
    userID, ok := userIDVal.(int)
    if !ok {
        http.Error(w, `{"error": "Invalid user ID"}`, http.StatusInternalServerError)
        return
    }


	var req ad.CreateAdRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := validateAd(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	query :=`INSERT INTO ads (title, description, image_url, price, user_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`

	var adID int
	var createdAt time.Time
	err = database.DB.QueryRow(
		query,
		req.Title,
		req.Description,
		req.ImageURL,
		req.Price,
		userID,
	).Scan(&adID, &createdAt)

	if err != nil {
		log.Printf("Failed to insert advertisement: %v", err)
		http.Error(w, "Failed to create advertisement", http.StatusInternalServerError)
		return
	}

	ad := ad.Ad{
		ID:          adID,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		UserID:      userID,
		CreatedAt:	 createdAt,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ad)
}

func validateAd(req ad.CreateAdRequest) error {
	if len(req.Title) < 3 || len(req.Title) > 100 {
		return errors.New("Title must be between 3 and 100 characters")
	}

	if len(req.Description) < 10 || len(req.Description) > 1000 {
		return errors.New("Description must be between 10 and 1000 characters")
	}


	if req.Price <= 0 {
		return errors.New("Price must be positive")
	}

	if !isValidImageURL(req.ImageURL) {
		return errors.New("Invalid image URL. Supported formats: jpg, jpeg, png, gif")
	}

	return nil
}

func isValidImageURL(url string) bool {

	if len(url) > 2048 {
		return false
	}


	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(strings.ToLower(url), ext) {
			return true
		}
	}
	return false
}


func GetAdsFeed(w http.ResponseWriter, r *http.Request) {
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }

    size, err := strconv.Atoi(r.URL.Query().Get("size"))
    if err != nil || size < 1 || size > 100 {
        size = 10
    }

    sortBy := r.URL.Query().Get("sort_by")
    if sortBy == "" {
        sortBy = "created_at"
    }

    sortOrder := r.URL.Query().Get("sort_order")
    if sortOrder == "" {
        sortOrder = "DESC"
    }

    minPrice, _ := strconv.ParseFloat(r.URL.Query().Get("min_price"), 64)
    maxPrice, _ := strconv.ParseFloat(r.URL.Query().Get("max_price"), 64)

    if minPrice < 0 {
        minPrice = 0
    }
    if maxPrice <= 0 {
        maxPrice = math.MaxFloat64
    }

    if minPrice > maxPrice {
        http.Error(w, "min_price cannot be greater than max_price", http.StatusBadRequest)
        return
    }

    allowedSortFields := map[string]bool{
        "created_at":   true,
        "price":        true,    
    }

    if !allowedSortFields[sortBy] {
        http.Error(w, "Invalid sort_by parameter", http.StatusBadRequest)
        return
    }

    if sortOrder != "ASC" && sortOrder != "DESC" {
        http.Error(w, "Invalid sort_order parameter", http.StatusBadRequest)
        return
    }

    var currentUserID int
    if userID := r.Context().Value("user_id"); userID != nil {
        currentUserID = userID.(int)
    }

    ads, total, err := getAdsFeed(page, size, sortBy, sortOrder, minPrice, maxPrice, currentUserID)
    if err != nil {
        log.Printf("Failed to get ads feed: %v", err)
        http.Error(w, "Failed to get ads feed", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "page":         page,
        "page_size":    size,
        "total_ads":    total,
        "total_pages":  int(math.Ceil(float64(total) / float64(size))),
        "ads":          ads,
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    json.NewEncoder(w).Encode(response)
}


func getAdsFeed(page, size int, sortBy, sortOrder string, minPrice, maxPrice float64, currentUserID int) ([]ad.FeedAd, int, error) {
	offset := (page-1) * size
	query := `
		SELECT 
			a.id, 
			a.title, 
			a.description, 
			a.image_url, 
			a.price, 
			v.login AS author_login,  
			a.created_at,             
			CASE WHEN a.user_id = $1 THEN true ELSE false END AS is_owner
		FROM ads a
		JOIN vkusers v ON a.user_id = v.id 
		WHERE a.price BETWEEN $2 AND $3
		ORDER BY ` + sortBy + ` ` + sortOrder + `
		LIMIT $4 OFFSET $5
	`

	countQuery := `
		SELECT COUNT(*)
		FROM	ads
		WHERE price BETWEEN $1 AND $2
	`

	var total int 
	err := database.DB.QueryRow(countQuery, minPrice, maxPrice).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	rows, err := database.DB.Query(query, currentUserID, minPrice, maxPrice, size, offset)
	if err != nil {
		return nil, 0, err
	}


	defer rows.Close()

	var ads []ad.FeedAd
	for rows.Next() {
		var ad ad.FeedAd

		err := rows.Scan(
			&ad.ID,
			&ad.Title,
			&ad.Description,
			&ad.ImageURL,
			&ad.Price,
			&ad.AuthorLogin,
			&ad.CreatedAt, 
			&ad.IsOwner,
		)	
		if err != nil {
			return nil, 0, err
		}
		ads = append(ads, ad)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return ads, total, nil

}




