aws iam create-role --role-name notification-cloudformation --assume-role-policy-document file://notification-cloudformation-assume-role-policy.json
aws iam put-role-policy --role-name notification-cloudformation --policy-name cloud-formation-role-policy --policy-document file://notification-cloudformation-role-policy.json

aws iam create-role --role-name notification-codepipeline --assume-role-policy-document file://notification-codepipeline-assume-role-policy.json
aws iam put-role-policy --role-name notification-codepipeline --policy-name cloud-formation-role-policy --policy-document file://notification-codepipeline-role-policy.json

aws iam create-role --role-name notification-codebuild --assume-role-policy-document file://notification-codebuild-assume-role-policy.json
aws iam put-role-policy --role-name notification-codebuild --policy-name cloud-formation-role-policy --policy-document file://notification-codebuild-role-policy.json 
