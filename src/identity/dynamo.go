package identity

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

// CreateEntry entity
func (i Identity)CreateEntry() (Identity, error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_DB_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_DB_ENDPOINT")),
	})
	if err != nil {
		return Identity{}, err
	}

	regs, err := convertRegistrationsToDynamo(i.Registrations)
	if err != nil {
		return Identity{}, err
	}

	svc := dynamodb.New(s)
	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("AWS_DB_TABLE")),
		Item: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(i.ID),
			},
			"email": {
				S: aws.String(i.Email),
			},
			"phone": {
				S: aws.String(i.Phone),
			},
			"company": {
				BOOL: aws.Bool(i.Company),
			},
			"registrations": &regs,
		},
		ConditionExpression: aws.String("attribute_not_exists(#IDENTIFIER)"),
		ExpressionAttributeNames: map[string]*string{
			"#IDENTIFIER": aws.String("identifier"),
		},
	}
	_, putErr := svc.PutItem(input)
	if putErr != nil {
		if awerr, ok := putErr.(awserr.Error); ok {
			switch awerr.Code(){
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return Identity{}, errors.New("identity already exists")
			}
		} else {
			return Identity{}, errors.New("unknown create error")
		}
	}

	return i, nil
}

// RetrieveEntry entity
func (i Identity)RetrieveEntry() (Identity, error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_DB_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_DB_ENDPOINT")),
	})
	if err != nil {
		return Identity{}, err
	}
	svc := dynamodb.New(s)
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(i.ID),
			},
		},
		TableName: aws.String(os.Getenv("AWS_DB_TABLE")),
	}
	result, err := svc.GetItem(input)
	if err != nil {
		return Identity{}, err
	}


	regLen := len(result.Item["registrations"].L)
	regs := []Registration{}
	if regLen >= 1 {
		for j := 0; j < regLen; j++ {
			ritem := result.Item["registrations"].L[j].M

			reg := Registration{
				Plate: *ritem["plate"].S,
				Oversized: *ritem["oversized"].BOOL,
				VehicleType: GetVehicleType(*ritem["vehicleType"].S),
			}

			regs = append(regs, reg)
		}
	}

	return Identity{
		ID: *result.Item["identifier"].S,
		Phone: *result.Item["phone"].S,
		Email: *result.Item["email"].S,
		Company: *result.Item["company"].BOOL,
		Registrations: regs,
	}, nil
}

// UpdateEntry entity
func (i Identity)UpdateEntry(n Identity) (Identity, error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_DB_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_DB_ENDPOINT")),
	})
	if err != nil {
		return Identity{}, err
	}

	regs, err := convertRegistrationsToDynamo(n.Registrations)
	if err != nil {
		return Identity{}, err
	}

	svc := dynamodb.New(s)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#EMAIL": aws.String("email"),
			"#PHONE": aws.String("phone"),
			"#COMPANY": aws.String("company"),
			"#REGISTRATIONS": aws.String("registrations"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(n.Email),
			},
			":phone": {
				S: aws.String(n.Phone),
			},
			":company": {
				BOOL: aws.Bool(n.Company),
			},
			":registrations": &regs,
		},
		Key: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(i.ID),
			},
		},
		TableName: aws.String(os.Getenv("AWS_DB_TABLE")),
		ReturnValues: aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET #EMAIL = :email, #PHONE = :phone, #COMPANY = :company, #REGISTRATIONS = :registrations"),
	}
	_, updateErr := svc.UpdateItem(input)
	if updateErr != nil {
		return Identity{}, updateErr
	}

	return n, nil
}

// DeleteEntry entity
func (i Identity)DeleteEntry() (Identity, error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_DB_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_DB_ENDPOINT")),
	})
	if err != nil {
		return Identity{}, nil
	}
	svc := dynamodb.New(s)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(i.ID),
			},
		},
		TableName: aws.String(os.Getenv("AWS_DB_TABLE")),
	}
	_, delErr := svc.DeleteItem(input)
	if delErr != nil {
		return Identity{}, delErr
	}

	return Identity{}, nil
}

// ScanEntry entry
func (i Identity)ScanEntry() (Identity, error) {
	if len(i.Registrations) == 0 {
		return Identity{}, errors.New("need at least 1 plate")
	}

	input := &dynamodb.ScanInput{}

	if i.Phone != "" {
		input = &dynamodb.ScanInput{
			ExpressionAttributeNames: map[string]*string{
				"#PHONE": aws.String("phone"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":phone": {
					S: aws.String(i.Phone),
				},
			},
			FilterExpression: aws.String("#PHONE = :phone"),
			TableName: aws.String(os.Getenv("AWS_DB_TABLE")),
		}
	} else if i.Email != "" {
		input = &dynamodb.ScanInput{
			ExpressionAttributeNames: map[string]*string{
				"#EMAIL": aws.String("email"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":email": {
					S: aws.String(i.Email),
				},
			},
			FilterExpression: aws.String("#EMAIL = :email"),
			TableName: aws.String(os.Getenv("AWS_DB_TABLE")),
		}
	} else if i.Phone != "" && i.Email != "" {
		input = &dynamodb.ScanInput{
			ExpressionAttributeNames: map[string]*string{
				"#EMAIL": aws.String("email"),
				"#PHONE": aws.String("phone"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":email": {
					S: aws.String(i.Email),
				},
				":phone": {
					S: aws.String(i.Phone),
				},
				":plate": {
					S: aws.String(i.Registrations[0].Plate),
				},
			},
			FilterExpression: aws.String("#PHONE = :phone AND #EMAIL = :email"),
			TableName: aws.String(os.Getenv("AWS_DB_TABLE")),
		}
	}


	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_DB_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_DB_ENDPOINT")),
	})
	if err != nil {
		return Identity{}, err
	}
	svc := dynamodb.New(s)
	result, err := svc.Scan(input)
	if err != nil {
		return Identity{}, err
	}

	if len(result.Items) == 1 {
		item := result.Items[0]

		ident := Identity{
			ID: *item["identifier"].S,
			Email: *item["email"].S,
			Phone: *item["phone"].S,
			Company: *item["company"].BOOL,
		}

		regLen := len(item["registrations"].L)
		if regLen >= 1 {
			regs := []Registration{}

			for j := 0; j < regLen; j++ {
				ritem := item["registrations"].L[j].M

				reg := Registration{
					Plate: *ritem["plate"].S,
					Oversized: *ritem["oversized"].BOOL,
					VehicleType: GetVehicleType(*ritem["vehicleType"].S),
				}

				regs = append(regs, reg)
			}

			ident.Registrations = regs
		}

		return ident, nil
	}

	return Identity{}, nil
}

// ScanEntries entities
func (i Identity)ScanEntries() ([]Identity, error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_DB_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_DB_ENDPOINT")),
	})
	if err != nil {
		return []Identity{}, err
	}
	svc := dynamodb.New(s)
	input := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("AWS_DB_TABLE")),
	}
	result, err := svc.Scan(input)
	if err != nil {
		return []Identity{}, err
	}

	itemLen := len(result.Items)
	if itemLen >= 1 {
		idents := []Identity{}

		for i := 0; i < itemLen; i++ {
			item := result.Items[i]

			ident := Identity{
				ID: *item["identifier"].S,
				Email: *item["email"].S,
				Phone: *item["phone"].S,
				Company: *item["company"].BOOL,
			}

			regLen := len(item["registrations"].L)
			if regLen >= 1 {
				regs := []Registration{}

				for j := 0; j < regLen; j++ {
					ritem := item["registrations"].L[j].M

					reg := Registration{
						Plate: *ritem["plate"].S,
						Oversized: *ritem["oversized"].BOOL,
						VehicleType: GetVehicleType(*ritem["vehicleType"].S),
					}

					regs = append(regs, reg)
				}

				ident.Registrations = regs
			}

			idents = append(idents, ident)
		}

		return idents, nil
	}

	return []Identity{}, nil
}

func convertRegistrationsToDynamo(regs []Registration) (dynamodb.AttributeValue, error) {
	retMap := map[string]*dynamodb.AttributeValue{}
	ret := dynamodb.AttributeValue{}

	if len(regs) >= 1 {
		for _, reg := range regs {
			retMap["vehicleType"] = &dynamodb.AttributeValue{
				S: aws.String(reg.VehicleType.convertToString()),
			}
			retMap["oversized"] = &dynamodb.AttributeValue{
				BOOL: aws.Bool(reg.Oversized),
			}
			retMap["plate"] = &dynamodb.AttributeValue{
				S: aws.String(reg.Plate),
			}
		}

		ret = dynamodb.AttributeValue{
			L: []*dynamodb.AttributeValue{
				{
					M: retMap,
				},
			},
		}
	}

	return ret, nil
}