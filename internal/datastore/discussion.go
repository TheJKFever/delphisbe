package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	logrus.Debug("GetDiscussionByID::Dynamo GetItem")
	res, err := d.dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.dbConfig.Discussions.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		logrus.WithError(err).Errorf("GetDiscussionByID: Failed getting discussion by ID (%s)", id)
		return nil, err
	}

	if res.Item == nil {
		return nil, nil
	}

	discussionObj := model.Discussion{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &discussionObj)

	if err != nil {
		logrus.WithError(err).Errorf("GetDiscussionsByID: Failed unmarshaling discussion by ID (%s)", id)
		return nil, err
	}

	return &discussionObj, nil
}

func (d *db) ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error) {
	logrus.Debug("ListDiscussions::Dynamo Scan")
	res, err := d.dynamo.Scan(&dynamodb.ScanInput{
		TableName: aws.String(d.dbConfig.Discussions.TableName),
	})

	if err != nil {
		logrus.WithError(err).Errorf("ListDiscussions: Failed listing discussions")
		return nil, err
	}

	if res.Count == nil || res.Items == nil {
		logrus.Errorf("ListDiscussions: Returned item set is nil")
	}

	ids := make([]string, 0)
	edges := make([]*model.DiscussionsEdge, 0)
	for _, elem := range res.Items {
		discussionObj := model.Discussion{}
		err := dynamodbattribute.UnmarshalMap(elem, &discussionObj)
		if err != nil {
			logrus.WithError(err).Warnf("ListDiscussion: Failed unmarshaling discussion: %+v", elem)
			continue
		}
		edges = append(edges, &model.DiscussionsEdge{
			Node: &discussionObj,
		})
		ids = append(ids, discussionObj.ID)
	}

	return &model.DiscussionsConnection{
		IDs:   ids,
		Edges: edges,
	}, nil
}

func (d *db) PutDiscussion(ctx context.Context, discussion model.Discussion) (*model.Discussion, error) {
	logrus.Debug("PutDiscussion::Dynamo PutItem")
	av, err := dynamodbattribute.MarshalMap(discussion)
	if err != nil {
		logrus.WithError(err).Errorf("PutDiscussion: Failed to marshal discussion object: %+v", discussion)
		return nil, err
	}
	_, err = d.dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.dbConfig.Discussions.TableName),
		Item:      av,
	})

	if err != nil {
		logrus.WithError(err).Errorf("PutDiscussion: Failed to put discussion object: %+v", av)
		return nil, err
	}
	return &discussion, nil
}
