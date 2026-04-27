package handlers

import (
	"at-backend-claims/internal/domain"
	"errors"
	"fmt"
	"unicode/utf8"
)

type claimWithoutCategory struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
}

func (claim claimWithoutCategory) Validate() error {
	var err error

	if claim.Title == nil {
		err = errors.Join(err, errors.New("title field is missing"))
	} else {
		if *claim.Title == "" {
			err = errors.Join(err, errors.New("title field is empty"))
		}

		if utf8.RuneCountInString(*claim.Title) > domain.ClaimTitleCharactersMaxCount {
			err = errors.Join(err, fmt.Errorf("title field value is too long; it can be <= %d characters", domain.ClaimTitleCharactersMaxCount))
		}
	}

	if claim.Description == nil {
		err = errors.Join(err, errors.New("description field is missing"))
	} else {
		if *claim.Description == "" {
			err = errors.Join(err, errors.New("description field is empty"))
		}

		if utf8.RuneCountInString(*claim.Description) > domain.ClaimDescriptionCharactersMaxCount {
			err = errors.Join(err, fmt.Errorf("description field value is too long; it can be <= %d characters", domain.ClaimDescriptionCharactersMaxCount))
		}
	}

	if claim.Latitude == nil {
		err = errors.Join(err, errors.New("latitude field is missing"))
	} else {
		if !(-90 <= *claim.Latitude || *claim.Latitude <= 90) {
			err = errors.Join(err, errors.New("latitude field has an invalid value"))
		}
	}

	if claim.Longitude == nil {
		err = errors.Join(err, errors.New("longitude field is missing"))
	} else {

		if !(-180 <= *claim.Longitude || *claim.Longitude <= 180) {
			err = errors.Join(err, errors.New("longitude field has an invalid value"))
		}
	}

	return err
}
