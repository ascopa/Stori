# Stori Software Engineer Technical Challenge

This repository presents a solution to the challenge presented in the **Stori Software Engineer Technical Challenge**.

---

## Architecture Overview

The solution consists of two main AWS Lambda functions:

- **`process-user-transaction`**
- **`send-user-email`**

Each Lambda is integrated using event-driven architecture with AWS EventBridge, S3, DynamoDB, and SES.

<img width="1031" height="578" alt="Challenge (1)" src="https://github.com/user-attachments/assets/6a85ca06-7339-4684-b222-687a0f3a7d7f" />

---

## Lambda: `process-user-transaction`

This function is responsible for retrieving a CSV file from the S3 bucket **`stori-user-transactions`** and processing the transactions it contains.

### How it works:

1. **Triggering**  
   The S3 bucket sends an event to the **default EventBridge bus**, which is captured by the rule **`user-transaction-to-process-transaction`**.  
   This rule extracts the bucket name and key from the event and triggers the Lambda.

2. **Processing**  
   The Lambda retrieves the file content and processes the transactions using Go routines and a worker pool (controlled by the `NUM_WORKERS` environment variable).  
   It performs two main actions:
    - Parses each CSV line into a transaction
    - Inserts each transaction into the DynamoDB table **`Transactions`**

3. **Performance**  
   Since both reading and insertion are fast operations, processing is parallelized in batches to improve performance.

4. **Event Emission**  
   Once processing is complete, the Lambda computes average statistics and emits a new event to the **default EventBridge bus**.

> Errors encountered during CSV line parsing or DB insertion are reported at the end but do **not interrupt** the process.
> 
> The use of EventBridge instead of direct S3 triggers avoids potential recursive execution issues.


---

## Lambda: `send-user-email`

This function is triggered by the event emitted after processing is complete.

### How it works:

1. **Triggering**  
   Captures events through the EventBridge rule **`user-transaction-to-email`**.

2. **Template Retrieval**  
   Loads an HTML email template from the S3 bucket **`stori-user-transactions-email-templates`**, allowing easy updates without redeploying the Lambda.

3. **User Lookup**  
   Uses the `AccountId` in the event to retrieve the user’s email from the **`Users`** DynamoDB table.

4. **Email Delivery**  
   Sends a transactional email via **Amazon SES** with the calculated statistics.

> This separation ensures that transaction processing and email sending are **loosely coupled and resilient**.

---

## Error Handling and DLQs

Both Lambdas are configured with **Dead Letter Queues (DLQs)** to:

- Avoid infinite error loops
- Enable post-mortem analysis of failed executions

---

## Testing Email Delivery (SES)

Since SES requires verified recipient addresses, a test user is preloaded in the **`Users`** table:

AccountId: AccountId123
Email: storichallenge@gmail.com
Name: John Doe

📧 **Inbox Credentials**:
- Email: `storichallenge@gmail.com`
- Password: `Challenge123*`

📁 **Important**: The test CSV file must include the AccountId `AccountId123` for email delivery to succeed.

---

## Testing

To trigger the process, a csv file must be uploaded to the **stori-user-transactions** bucket. The AWS IAM Test user is the following:

User name,Password,Console sign-in URL

StoriChallengeUser,)cd89O1(,https://010438498297.signin.aws.amazon.com/console

The test example file is in the root dir and its named transactions3000.csv 


## Future improvements
Its expected that each csv file contains the transactions belonging to only one user. There are no validations for multiple accountIds. The lambda could be refactored with this behaviour in mind.

Processing very large files is currently limited by the execution time of a Lambda function.  
For longer-running processes—such as handling transactions for many users simultaneously—it would be more appropriate to use a compute service like **AWS Fargate**.


