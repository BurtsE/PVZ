package receptions

import domain "pvz/internal/domain/receptions"

func convertReceptionsMapToDomain(m map[receptionDB][]productDb) ([]domain.Reception, error) {
	receptions := make([]domain.Reception, 0, len(m))
	for reception, products := range m {
		result := make([]domain.Product, 0, len(products))
		for _, p := range products {
			product, err := domain.NewProduct(p.id, p.arrivalTime, p.productType)
			if err != nil {
				return nil, err
			}
			result = append(result, *product)
		}
		r, err := domain.NewReception(reception.id, reception.status, reception.registrationDate, result)
		if err != nil {
			return nil, err
		}
		receptions = append(receptions, *r)
	}
	return receptions, nil
}
