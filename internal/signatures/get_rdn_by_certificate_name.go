package signatures

import (
	DB "couriergate/internal/db"
	"errors"
	"fmt"
)

// ----- получить КПС по имени сертификата -----
// 01
func GetRDNByCertificateName(certificateName string) (string, error) {
	var RDN string

	rows, err := DB.DB_COURIER.Query("SELECT TOP 1 RDN FROM RDN (NOLOCK) WHERE CertificateName = @p1", certificateName)
	if err != nil {
		return RDN, errors.New("Error (SU-040101): " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&RDN)
		if err != nil {
			return RDN, errors.New("Error (SU-040102): " + fmt.Sprintf("%s\n", err))
		}
		return RDN, nil
	}

	return "", errors.New("Error (SU-040103): неизвестен RDN для сертификата " + fmt.Sprintf("%s\n", certificateName))
}
