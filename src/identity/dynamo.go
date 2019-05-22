package identity

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

// CreateEntry entity
func (i Identity) CreateEntry() (Identity, error) {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("AWS_DB_REGION")),
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
		ReturnValues: aws.String("ALL_OLD"),
	}
	_, putErr := svc.PutItem(input)
	if putErr != nil {
		if awsErr, ok := putErr.(awserr.Error); ok {
			switch awsErr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return Identity{}, errors.New("identity already exists")
			default:
				return Identity{}, fmt.Errorf("unknown err: %v", awsErr)
			}
		} else {
			return Identity{}, fmt.Errorf("unknown create error err: %v", putErr)
		}
	}

	// fmt.Println(fmt.Sprintf("create result: %v", res))

	return i, nil
}

// RetrieveEntry entity
func (i Identity) RetrieveEntry() (Identity, error) {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("AWS_DB_REGION")),
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

	return convertDynamoToIdentity(result.Item)
}

// UpdateEntry entity
func (i Identity) UpdateEntry(n Identity) (Identity, error) {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("AWS_DB_REGION")),
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

	input := &dynamodb.UpdateItemInput{}
	if n.Email != "" {
		input = &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#EMAIL": aws.String("email"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":email": {
					S: aws.String(n.Email),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"identifier": {
					S: aws.String(i.ID),
				},
			},
			TableName:        aws.String(os.Getenv("AWS_DB_TABLE")),
			ReturnValues:     aws.String("ALL_NEW"),
			UpdateExpression: aws.String("SET #EMAIL = :email"),
		}
	} else if n.Phone != "" {
		input = &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#PHONE": aws.String("phone"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":phone": {
					S: aws.String(n.Phone),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"identifier": {
					S: aws.String(i.ID),
				},
			},
			TableName:        aws.String(os.Getenv("AWS_DB_TABLE")),
			ReturnValues:     aws.String("ALL_NEW"),
			UpdateExpression: aws.String("SET #PHONE = :phone"),
		}
	} else if len(n.Registrations) >= 1 {
		input = &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#REGISTRATIONS": aws.String("registrations"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":registrations": &regs,
			},
			Key: map[string]*dynamodb.AttributeValue{
				"identifier": {
					S: aws.String(i.ID),
				},
			},
			TableName:        aws.String(os.Getenv("AWS_DB_TABLE")),
			ReturnValues:     aws.String("ALL_NEW"),
			UpdateExpression: aws.String("SET #REGISTRATIONS = :registrations"),
		}
	} else if i.Company != n.Company {
		input = &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#COMPANY": aws.String("company"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":company": {
					BOOL: aws.Bool(n.Company),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"identifier": {
					S: aws.String(i.ID),
				},
			},
			TableName:        aws.String(os.Getenv("AWS_DB_TABLE")),
			ReturnValues:     aws.String("ALL_NEW"),
			UpdateExpression: aws.String("SET #COMPANY = :company"),
		}
	}

	ret, err := svc.UpdateItem(input)
	if err != nil {
		return Identity{}, err
	}

	return convertDynamoToIdentity(ret.Attributes)
}

// DeleteEntry entity
func (i Identity) DeleteEntry() (Identity, error) {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("AWS_DB_REGION")),
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
func (i Identity) ScanEntry() (Identity, error) {
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
			TableName:        aws.String(os.Getenv("AWS_DB_TABLE")),
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
			TableName:        aws.String(os.Getenv("AWS_DB_TABLE")),
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
			},
			FilterExpression: aws.String("#PHONE = :phone AND #EMAIL = :email"),
			TableName:        aws.String(os.Getenv("AWS_DB_TABLE")),
		}
	}

	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("AWS_DB_REGION")),
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
		ident, err := convertDynamoToIdentity(item)
		if err != nil {
			return Identity{}, errors.New("could't return ident")
		}

		return ident, nil
	}

	return Identity{}, nil
}

// ScanAll scan all the entries
func ScanAll() ([]Identity, error) {
	i := Identity{}
	return i.ScanEntries()
}

// ScanEntries entities
func (i Identity) ScanEntries() ([]Identity, error) {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("AWS_DB_REGION")),
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
			ident, err := convertDynamoToIdentity(item)
			if err != nil {
				return idents, errors.New("couldn't return idents")
			}
			idents = append(idents, ident)
		}

		return idents, nil
	}

	return []Identity{}, nil
}

func convertRegistrationsToDynamo(regs []Registration) (dynamodb.AttributeValue, error) {
	ret := dynamodb.AttributeValue{}
	lMap := []*dynamodb.AttributeValue{}

	if len(regs) >= 1 {
		for _, reg := range regs {
			retMap := map[string]*dynamodb.AttributeValue{}
			retMap["vehicleType"] = &dynamodb.AttributeValue{
				S: aws.String(reg.VehicleType.convertToString()),
			}
			retMap["oversized"] = &dynamodb.AttributeValue{
				BOOL: aws.Bool(reg.Oversized),
			}
			retMap["plate"] = &dynamodb.AttributeValue{
				S: aws.String(reg.Plate),
			}
			mmap := &dynamodb.AttributeValue{
				M: retMap,
			}


			lMap = append(lMap, mmap)
		}

		ret = dynamodb.AttributeValue{
			L: lMap,
		}
	} else {
		ret = dynamodb.AttributeValue{
			BOOL: aws.Bool(false),
		}
	}

	return ret, nil
}

func convertDynamoToRegistration(reg *dynamodb.AttributeValue) (Registration, error) {
	ret := Registration{}

	for key, value := range reg.M {
		switch key {
		case "plate":
			ret.Plate = *value.S
		case "vehicleType":
			ret.VehicleType = GetVehicleType(*value.S)
		case "oversized":
			ret.Oversized = *value.BOOL
		}
	}
	if ret.Plate != "" {
		return ret, nil
	}

	return Registration{}, errors.New("couldn't convert to registration")
}

func convertDynamoToRegistrations(items []*dynamodb.AttributeValue) ([]Registration, error) {
	ret := []Registration{}

	for _, item := range items {
		reg, err := convertDynamoToRegistration(item)
		if err != nil {
			return []Registration{}, err
		}

		ret = append(ret, reg)
	}

	return ret, nil
}

func convertDynamoToIdentity(items map[string]*dynamodb.AttributeValue) (Identity, error) {
	ret := Identity{}
	for key, value := range items {
		switch key {
		case "registrations":
			regs, err := convertDynamoToRegistrations(value.L)
			if err != nil {
				return Identity{}, err
			}
			ret.Registrations = regs

		case "company":
			ret.Company = *value.BOOL

		case "phone":
			ret.Phone = *value.S

		case "email":
			ret.Email = *value.S

		case "identifier":
			ret.ID = *value.S
		}
	}

	if ret.ID != "" {
		return ret, nil
	}

	return Identity{}, errors.New("couldn't convert to identity")
}
