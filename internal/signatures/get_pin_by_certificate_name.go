package signatures

import (
	DB "couriergate/internal/db"
	"errors"
	"fmt"
)

// ----- получить Pin по имени сертификата -----
// 01
func GetPinByCertificateName(certificateName string) (string, error) {
	var Pin string

	rows, err := DB.DB_COURIER.Query("SELECT TOP 1 Pin FROM PIN (NOLOCK) WHERE CertificateName = @p1", certificateName)
	if err != nil {
		return Pin, errors.New("Error (SU-050101): " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&Pin)
		if err != nil {
			return Pin, errors.New("Error (SU-050102): " + fmt.Sprintf("%s\n", err))
		}
		return Pin, nil
	}

	return "", errors.New("Error (SU-050103): неизвестен Pin для сертификата " + fmt.Sprintf("%s\n", certificateName))
}
