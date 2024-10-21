package cmd

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"

    "github.com/olekukonko/tablewriter"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
)

var filterCmd = &cobra.Command{
    Use:   "filter",
    Short: "Filter trails based on criteria",
    Long:  `Filter trails by various criteria such as restrooms, fishing, bike, horse, fee, recycle bin, grills, bike rack, and dog tube, with pagination support.`,
    Run:   filterTrails,
}

func init() {
    filterCmd.Flags().String("restrooms", "", "Filter by restrooms (Yes/No)")
    filterCmd.Flags().String("fishing", "", "Filter by fishing (Yes/No)")
    filterCmd.Flags().String("bike", "", "Filter by bike trail (Yes/No)")
    filterCmd.Flags().String("horse", "", "Filter by horse trail (Possible/Not Recommended/Designated/NA)")
    filterCmd.Flags().String("fee", "", "Filter by fee (Yes/No)")
    filterCmd.Flags().String("recycle_bin", "", "Filter by recycle bin (Yes/No)")
    filterCmd.Flags().String("grills", "", "Filter by grills (Yes/No)")
    filterCmd.Flags().String("bike_rack", "", "Filter by bike rack (Yes/No)")
    filterCmd.Flags().String("dog_tube", "", "Filter by dog tube (Yes/No)")
    filterCmd.Flags().Int("page", 1, "Page number for pagination")
    filterCmd.Flags().Int("limit", 10, "Number of results per page for pagination")

    rootCmd.AddCommand(filterCmd)
}

func filterTrails(cmd *cobra.Command, args []string) {
    // Get filter flags
    restrooms, _ := cmd.Flags().GetString("restrooms")
    fishing, _ := cmd.Flags().GetString("fishing")
    bikeTrail, _ := cmd.Flags().GetString("bike")
    horseTrail, _ := cmd.Flags().GetString("horse")
    fee, _ := cmd.Flags().GetString("fee")
    recycleBin, _ := cmd.Flags().GetString("recycle_bin")
    grills, _ := cmd.Flags().GetString("grills")
    bikeRack, _ := cmd.Flags().GetString("bike_rack")
    dogTube, _ := cmd.Flags().GetString("dog_tube")
    page, _ := cmd.Flags().GetInt("page")
    limit, _ := cmd.Flags().GetInt("limit")

    // Convert filters to lowercase for consistency
    filters := []string{}
    if restrooms != "" {
        filters = append(filters, "restrooms="+strings.ToLower(restrooms))
    }
    if fishing != "" {
        filters = append(filters, "fishing="+strings.ToLower(fishing))
    }
    if bikeTrail != "" {
        filters = append(filters, "bike_trail="+strings.ToLower(bikeTrail))
    }
    if horseTrail != "" {
        filters = append(filters, "horse_trail="+strings.ToLower(horseTrail))
    }
    if fee != "" {
        filters = append(filters, "fee="+strings.ToLower(fee))
    }
    if recycleBin != "" {
        filters = append(filters, "recycle_bin="+strings.ToLower(recycleBin))
    }
    if grills != "" {
        filters = append(filters, "grills="+strings.ToLower(grills))
    }
    if bikeRack != "" {
        filters = append(filters, "bike_rack="+strings.ToLower(bikeRack))
    }
    if dogTube != "" {
        filters = append(filters, "dog_tube="+strings.ToLower(dogTube))
    }

    // Add pagination parameters
    filters = append(filters, fmt.Sprintf("page=%d", page))
    filters = append(filters, fmt.Sprintf("limit=%d", limit))

    apiURL := "http://localhost:8080/trails?" + strings.Join(filters, "&")

    // Make the API request
    resp, err := http.Get(apiURL)
    if err != nil {
        logrus.Errorf("Error fetching trails: %v", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        logrus.Errorf("Failed to fetch trails: %s", resp.Status)
        return
    }

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logrus.Errorf("Error reading response: %v", err)
        return
    }

    // Parse the JSON response
    var response struct {
        Page    int             `json:"page"`
        Limit   int             `json:"limit"`
        Results []TrailResponse `json:"results"`
    }
    err = json.Unmarshal(body, &response)
    if err != nil {
        logrus.Errorf("Error parsing JSON response: %v", err)
        return
    }

    // Display the filtered results in a table format
    if len(response.Results) == 0 {
        logrus.Info("No trails found for the given criteria.")
        return
    }

    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Name", "Restrooms", "Picnic", "Fishing", "Difficulty", "Access Type", "TH Leash", "Bike Trail", "Horse Trail", "Fee", "Recycle Bin", "Grills", "Bike Rack", "Dog Tube"})

    for _, trail := range response.Results {
        table.Append([]string{
            trail.Name, trail.Restrooms, trail.Picnic, trail.Fishing, trail.Difficulty, trail.AccessType,
            trail.THLeash, trail.BikeTrail, trail.HorseTrail, trail.Fee, trail.RecycleBin, trail.Grills, trail.BikeRack, trail.DogTube,
        })
    }

    logrus.Infof("Showing page %d with %d results per page:", response.Page, response.Limit)
    table.Render()
}

// TrailResponse represents the structure of a trail in the response, excluding FID
type TrailResponse struct {
    Name        string `json:"name"`
    Restrooms   string `json:"restrooms"`
    Picnic      string `json:"picnic"`
    Fishing     string `json:"fishing"`
    Difficulty  string `json:"difficulty"`
    AccessType  string `json:"access_type"`
    THLeash     string `json:"th_leash"`
    BikeTrail   string `json:"bike_trail"`
    HorseTrail  string `json:"horse_trail"`
    Fee         string `json:"fee"`
    RecycleBin  string `json:"recycle_bin"`
    Grills      string `json:"grills"`
    BikeRack    string `json:"bike_rack"`
    DogTube     string `json:"dog_tube"`
}
