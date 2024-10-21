package handlers

import (
    "context"
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "strings"
    "trail-finder/models"
    "github.com/sirupsen/logrus"
    "trail-finder/db"
)

// LoadTrailsFromRequest handles loading a new CSV file from a client request
func LoadTrailsFromRequest(w http.ResponseWriter, r *http.Request) {
    var request struct {
        FilePath string `json:"file_path"`
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        logrus.Error("Invalid request body")
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := json.Unmarshal(body, &request); err != nil {
        logrus.Error("Invalid JSON")
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    filePath := request.FilePath
    if filePath == "" {
        logrus.Error("File path must be provided")
        http.Error(w, "File path must be provided", http.StatusBadRequest)
        return
    }

    if err := LoadTrails(filePath); err != nil {
        logrus.Errorf("Error loading trails: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    logrus.Infof("Trails loaded successfully from: %s", filePath)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Trails loaded successfully from: " + filePath))
}

// LoadDefaultData loads default data from a CSV if the trails table is empty
func LoadDefaultData(filename string) error {
    isEmpty, err := IsTableEmpty()
    if err != nil {
        logrus.Errorf("Failed to check if table is empty: %v", err)
        return err
    }

    if isEmpty {
        logrus.Infof("Loading default data from: %s", filename)
        return LoadTrails(filename)
    }

    logrus.Info("Trails data already exists. No default data loaded.")
    return nil
}

// IsTableEmpty checks if the trails table is empty
func IsTableEmpty() (bool, error) {
    if db.DbConn == nil {
        return false, fmt.Errorf("database connection is not initialized")
    }

    var count int
    err := db.DbConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM trails").Scan(&count)
    if err != nil {
        logrus.Errorf("Failed to check table row count: %v", err)
        return false, err
    }

    return count == 0, nil
}


// LoadTrails loads data from CSV into PostgreSQL, replacing existing data
func LoadTrails(filename string) error {
    // Check if the database connection is initialized
    if db.DbConn == nil {
        return fmt.Errorf("database connection is not initialized")
    }

    file, err := os.Open(filename)
    if err != nil {
        logrus.Errorf("Could not open file: %v", err)
        return fmt.Errorf("could not open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    rows, err := reader.ReadAll()
    if err != nil {
        logrus.Errorf("Could not read CSV: %v", err)
        return fmt.Errorf("could not read CSV: %w", err)
    }

    // Start a transaction to ensure atomicity
    tx, err := db.DbConn.Begin(context.Background())
    if err != nil {
        logrus.Errorf("Could not start transaction: %v", err)
        return fmt.Errorf("could not start transaction: %w", err)
    }

    // Clear the existing data in the trails table
    _, err = tx.Exec(context.Background(), "DELETE FROM trails")
    if err != nil {
        tx.Rollback(context.Background())
        logrus.Errorf("Could not clear existing data: %v", err)
        return fmt.Errorf("could not clear existing data: %w", err)
    }

    // Insert new data from CSV
    for i, row := range rows {
        if i == 0 || len(row) < 32 {
            continue // Skip header or invalid rows
        }

        fid, err := strconv.Atoi(row[0])
        if err != nil {
            continue // Skip rows where FID is not an integer
        }

        name := strings.ToLower(row[30])
        restrooms := strings.ToLower(row[1])
        picnic := strings.ToLower(row[2])
        fishing := strings.ToLower(row[3])
        bikeTrail := strings.ToLower(row[11])
        horseTrail := strings.ToLower(row[25])
        difficulty := strings.ToLower(row[21])
        fee := strings.ToLower(row[10])
        recycleBin := strings.ToLower(row[28])
        grills := strings.ToLower(row[13])
        bikeRack := strings.ToLower(row[11])
        dogTube := strings.ToLower(row[12])

        _, err = tx.Exec(context.Background(), `
            INSERT INTO trails (fid, name, restrooms, picnic, fishing, type, difficulty, access_type, th_leash, bike_trail, horse_trail, fee, recycle_bin, grills, bike_rack, dog_tube )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
        `, fid, name, restrooms, picnic, fishing, row[8], difficulty, row[6], row[31], bikeTrail, horseTrail, fee, recycleBin, grills, bikeRack, dogTube)

        if err != nil {
            tx.Rollback(context.Background())
            logrus.Errorf("Failed to insert data: %v", err)
            return fmt.Errorf("failed to insert data: %w", err)
        }
    }

    err = tx.Commit(context.Background())
    if err != nil {
        logrus.Errorf("Could not commit transaction: %v", err)
        return fmt.Errorf("could not commit transaction: %w", err)
    }

    logrus.Infof("Trails data replaced successfully from: %s", filename)
    return nil
}

// GetTrails handles GET requests to filter trails from PostgreSQL
func GetTrails(w http.ResponseWriter, r *http.Request) {
    // Fetch filter query parameters
    restrooms := strings.ToLower(r.URL.Query().Get("restrooms"))
    fishing := strings.ToLower(r.URL.Query().Get("fishing"))
    bikeTrail := strings.ToLower(r.URL.Query().Get("bike_trail"))
    horseTrail := strings.ToLower(r.URL.Query().Get("horse_trail"))

    // Pagination parameters
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    offset := (page - 1) * limit

    // Initialize query and args
    query := "SELECT * FROM trails WHERE 1=1"
    args := []interface{}{}
    i := 1

    // Build the query based on the filters
    if restrooms != "" {
        query += fmt.Sprintf(" AND LOWER(restrooms) = $%d", i)
        args = append(args, restrooms)
        i++
    }
    if fishing != "" {
        query += fmt.Sprintf(" AND LOWER(fishing) = $%d", i)
        args = append(args, fishing)
        i++
    }
    if bikeTrail != "" {
        query += fmt.Sprintf(" AND LOWER(bike_trail) = $%d", i)
        args = append(args, bikeTrail)
        i++
    }
    if horseTrail != "" {
        query += fmt.Sprintf(" AND LOWER(horse_trail) = $%d", i)
        args = append(args, horseTrail)
        i++
    }

    // Add pagination to the query
    query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
    args = append(args, limit, offset)

    logrus.Infof("SQL Query: %s", query)
    logrus.Infof("Arguments: %v", args)

    // Execute the query
    rows, err := db.DbConn.Query(context.Background(), query, args...)
    if err != nil {
        logrus.Errorf("Failed to query trails: %v", err)
        http.Error(w, "Failed to query trails", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Prepare the response
    var filteredTrails []models.Trail
    for rows.Next() {
        var trail models.Trail
        err := rows.Scan(&trail.FID, &trail.Name, &trail.Restrooms, &trail.Picnic, &trail.Fishing, &trail.Type, &trail.Difficulty, &trail.AccessType, &trail.THLeash, &trail.BikeTrail, &trail.HorseTrail, &trail.Fee, &trail.RecycleBin, &trail.Grills, &trail.BikeRack, &trail.DogTube)
        if err != nil {
            logrus.Errorf("Failed to scan trails: %v", err)
            http.Error(w, "Failed to scan trails", http.StatusInternalServerError)
            return
        }
        filteredTrails = append(filteredTrails, trail)
    }

    // Set the response header to JSON
    w.Header().Set("Content-Type", "application/json")

    // Prepare the response with pagination info
    response := map[string]interface{}{
        "page":    page,
        "limit":   limit,
        "results": filteredTrails,
    }

    // Respond with the filtered trails
    logrus.Infof("Responding with %d results for page %d", len(filteredTrails), page)
    json.NewEncoder(w).Encode(response)
}