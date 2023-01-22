package cashiers

import (
	"encoding/hex"
	"regexp"

	"github.com/tamboto2000/gotoko-pos/apperror"
	"github.com/tamboto2000/gotoko-pos/helpers/aescrypt"
	"github.com/tamboto2000/gotoko-pos/models"
)

type Cashiers struct {
	Cashiers []Cashier    `json:"cashiers"`
	Meta     CashiersMeta `json:"meta"`
}

type CashiersMeta struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}

type Cashier struct {
	Sessions []models.CashierSession `json:"-"`
	models.Cashier
	encryptPasscode string
}

var (
	ErrNameEmpty   = apperror.New(`"name" is required`, apperror.AnyRequired, "name")
	ErrNameInvalid = apperror.New(
		`"name" is invalid, valid input is alphabets (A-Z) and numeric (0-9) with 4 minimum characters and 100 maximum characters`,
		apperror.AnyInvalid,
		"name",
	)
	ErrPasscodeEmpty          = apperror.New(`"passcode" is required`, apperror.AnyRequired, "passcode")
	ErrPasscodeInvalid        = apperror.New(`"passcode" is invalid, valid input is 6 numeric (0-9) charaters`, apperror.AnyInvalid, "passcode")
	ErrCashierNotFound        = apperror.New(`"Cashier not found"`, apperror.NotFound, "")
	ErrCashierSessionNotFound = apperror.New("Cashier session not found", apperror.NotFound, "")
)

func NewCashier(name, passcode string) (Cashier, error) {
	cashier := Cashier{
		Cashier: models.Cashier{
			Name:     name,
			Passcode: passcode,
		},
	}

	if err := cashier.Validate(); err != nil {
		return Cashier{}, err
	}

	return cashier, nil
}

func (c *Cashier) Validate() error {
	errList := apperror.NewErrorList()
	errList.SetType(apperror.AnyRequired)
	errList.SetPrefix("body ValidationError: ")

	// validate name
	if c.Name == "" {
		errList.Add(ErrNameEmpty)
	} else {
		rgx := regexp.MustCompile(`^[0-9A-Za-z ]{4,100}$`)
		if ok := rgx.MatchString(c.Name); !ok {
			errList.Add(ErrNameInvalid)
		}
	}

	// validate passcode
	if c.Passcode == "" {
		errList.Add(ErrPasscodeEmpty)
	} else {
		rgx := regexp.MustCompile(`^[0-9]{6,6}$`)
		if ok := rgx.MatchString(c.Passcode); !ok {
			errList.Add(ErrPasscodeInvalid)
		}
	}

	return errList.BuildError()
}

func (c *Cashier) ValidateForUpdate() error {
	errList := apperror.NewErrorList()
	errList.SetPrefix("body ValidationError: ")
	errList.SetType(apperror.ObjecMissing)

	if c.Name == "" && c.Passcode == "" {
		errList.Add(apperror.NewWithPeers(
			`"value" must contain at least one of [name]`,
			apperror.ObjecMissing,
			[]string{},
			"value",
			[]string{"name"},
		))
	}

	if c.Passcode != "" {
		rgx := regexp.MustCompile(`^[0-9]{6,6}$`)
		if ok := rgx.MatchString(c.Passcode); !ok {
			errList.Add(ErrPasscodeInvalid)
		}
	}

	return errList.BuildError()
}

func (c *Cashier) ValidateForLogin() error {
	errList := apperror.NewErrorList()
	errList.SetPrefix("body ValidationError: ")
	errList.SetType(apperror.AnyRequired)

	if c.Passcode == "" {
		errList.Add(ErrPasscodeEmpty)
	}

	return errList.BuildError()
}

func (c *Cashier) SetName(name string) {
	c.Cashier.Name = name
}

func (c *Cashier) SetPasscode(passcode string) {
	c.Cashier.Passcode = passcode
}

// EncryptPasscode encrypt passcode with AES encryption
func (c *Cashier) EncryptPasscode() error {
	enc, err := aescrypt.Encrypt([]byte(c.Passcode))
	if err != nil {
		return err
	}

	c.encryptPasscode = hex.EncodeToString(enc)
	return nil
}

func (c *Cashier) DecryptPasscode() error {
	raw, err := hex.DecodeString(c.encryptPasscode)
	if err != nil {
		return err
	}

	passcode, err := aescrypt.Decrypt(raw)
	if err != nil {
		return err
	}

	c.Passcode = string(passcode)
	return nil
}

func (c *Cashier) GetEncryptPasscode() string {
	return c.encryptPasscode
}

func (c *Cashier) GetName() string {
	return c.Cashier.Name
}

func (c *Cashier) GetPasscode() string {
	return c.Cashier.Passcode
}

func (c *Cashier) GetId() int {
	return c.Cashier.Id
}
