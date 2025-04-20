# Sprint 3  
## Frontend Development

## User stories
- Populating real time news in the Dashboard
- Frontend and backend integration for Start a post functionalities
- Layout for View Profile
- Cypress and unit testing
- Search feature for searching people 
- Designed layout for user profile pages
- Network feature to display my connections

## What issues your team planned to address
- Modify Cypress and unit tests completed in Sprint 2, introduction of new unit tests for : rendering all components in the dashboard, Start a post functionality, View Profile.
- As an authenticated user, render real time news on the left side of the dashboard and upon clicking the link redirects the user to that website.
- As an authenticated user, should be able to create a post using Start a post button, view all posts, edit and delete a particular post when Edit or Delete button is clicked. Render a form to allow a user to create a post.
    - Integrate front end and back end for Start a post functionality.
- As an authenticated user, should be able to view the profile when View Profile button is clicked.
- Search feature to search for users and user profiles on the platform
- As a user, I should be able to send connection request to user, accept and reject user connection requests
- As a user, I should be able to view users in my connection

## Which ones were successfully completed
- Successfully tested whether the Start a post button was rendered, the form was rendered when Start a post button was clicked, profile was rendered when View Profile button was clicked, creation, update, delete and display of posts.
- Populated real time news on the left side of the dashboard when an authenticated user is redireccted to the dashboard using the axios module.
- Render a form when Start a post button is clicked and enter the details in order to create a post.
- Integrated front and back end for create, update, display and deletion of posts.
- Designed basic layout for View Profile and render the page when an authenticated user clicks View Profile and display the Recent Activity.
- Designed basic layout for user profiles and render user profiles having about, recent activity and experience with connection request
- Created search bar, and dynamic user profile pages to display searched profiles
- Created network feature to display users in my connection

Details fetched when creating a post:
```

{
"enterpreneur_email": "johndoe@gmail.com",
"startup_name": "Smart AI Systems",
"industry" : "AI",
"description": "Building AI-powered tools for businesses",
"budget": 750000,
"timeframe": "8 months",
}
```

Details fetched when editing a post:
```

{
"startup_id":"1",
"email": "johndoe@gmail.com",
"startup_name": "Smart AI Systems",
"industry" : "AI",
"description": "Building AI-powered tools for businesses",
"budget": 850000,
"timeframe": "10 months",
}
```

Details fetched when deleting a post:
```

{
"email": "johndoe@gmail.com",
"startup_id":"1",
}
```

- Integrate front and Backend using API queries.
    - Create a post using a POST method
    - Display all posts usinf GET method
    - Edit a post using PUT method
    - Delete a post usig DELETE method

## Frontend Unit Tests

### Technologies Used
- Jest

### Test Cases
1. Dashboard Testing
    1. Render start a post button after user is authenticated.
    2. Render the start a post form when Start a post button is clciked.
    3. Render View Profile when user clicks View Profile button.
    4. Handle post creation.
    5. Handle post display.
    6. Handle post deletion.
    7. Handle post updation.

2. View Profile Testing
    1. Displays user information after authentication : check if details such as Name, Role, Location, no of connection, Upload Profile Picture button, About, Recent Activity and Experience is displayed properly.

### Cypress tests
- Modified the cypress test completed in sprint 2.



## Backend Development
## API Documentation

## User Stories Implemented
- As a user, I want to like a startup so I can show interest in it
- As a user, I want to comment on startups to share my thoughts
- As a user, I want to report a startup if I find inappropriate content
- As a user, I want to search for other users by name to find and connect with them
- As a platform, I want to enable universal connection requests between all roles
- As a platform, I want to implement intelligent connection suggestions based on user roles and profiles
- As a platform, I want to display startups in sorted order (newest first)

## What ones were successfully completed 
- Implemented like/unlike endpoints with duplicate prevention
- Built comment creation and retrieval system with validation
- Developed reporting system with duplicate prevention
- Created case-insensitive partial name search functionality
- Enabled universal connections between all user roles
- Implemented suggestion algorithm with role-based matching (100/80/50 scoring)
- Updated all startup queries to sort by created_at/updated_at
- Added comprehensive unit tests for all new endpoints
- Verified search and suggestion logic through testing


## User Management
1. Signup
- Endpoint: POST /v1/signup
- Description: Register a new user
- Request Body:
```json
{
  "full_name": "string",
  "email": "string",
  "role": "string (Entrepreneur/Investor/Mentor)",
  "phone_number": "string",
  "location": "string",
  "linkedin_profile": "string",
  "verification_proof": "string (optional)",
  
   "Entrepreneur-specific fields"
  "startup_name": "string (optional)",
  "industry": "string (optional)",
  "description": "string (optional)",
  "budget": "number (optional)",
  "timeframe": "string (optional)",
  
   "Investor-specific fields"
  "investment_range": "string (optional)",
  "preferred_industries": "string (optional)",
  "experience_years": "number (optional)",
  
   "Mentor-specific fields"
  "expertise": "string (optional)",
  "past_mentorships": "string (optional)",
  "years_of_experience": "number (optional)"
}
```

- Success Response:

```json
{
  "status": "success",
  "message": "User registered successfully! Please verify your account.",
  "user_id": "number"
}
```

2. Check User Existence
- Endpoint: GET /v1/user/check?email=user@example.com
- Description: Check if a user exists
- Query Parameters:
    email (required)
- Success Response:
```json
{
  "exists": true,
  "role": "string",
  "full_name": "string",
  "email": "string",
  "phone_number": "string",
  "location": "string",
  "linkedin": "string",
  "verification_status": "string"
}
```

3. Search Users by Name
- Endpoint: GET /v1/user/search?name=searchTerm
- Description: Search users by name (case-insensitive partial match)
- Query Parameters:
name (required)

Success Response:
```json
{
  "status": "success",
  "count": "number",
  "users": [
    {
      "id": "number",
      "full_name": "string",
      "email": "string",
      "role": "string",
      "location": "string",
      "linkedin_profile": "string"
    }
  ]
}
```

## Connection Management

1. Request Connection
- Endpoint: POST /v1/connection/request
- Description: Send a connection request
- Request Body:
```json
{
  "sender_email": "string",
  "receiver_email": "string"
}
```
- Success Response:
```json
{
  "status": "success",
  "message": "Connection request sent from sender@example.com to receiver@example.com"
}
```

2. Get Pending Connections
- Endpoint: GET /v1/connection/pending?email=user@example.com
- Description: Get pending connection requests for a user
- Query Parameters:
email (required)

Success Response:
```
[
  {
    "connection_id": number,
    "full_name": "string",
    "email": "string",
    "role": "string",
    "location": "string",
    "linkedin": "string"
  }
]
```

3. Accept Connection
- Endpoint: PUT /v1/connection/accept/{id}?email=user@example.com
- Description: Accept a pending connection request
- URL Parameters:
id (connection ID)

- Query Parameters:
email (required)

- Success Response:
```json
{
  "status": "success",
  "message": "Connection request [id] accepted"
}
```

4. Reject Connection
- Endpoint: DELETE /v1/connection/reject/{id}?email=user@example.com
- Description: Reject a pending connection request
- URL Parameters:
id (connection ID)

- Query Parameters:
email (required)

Success Response:
```json
{
  "status": "success",
  "message": "Connection request [id] rejected"
}
```

5. Get Accepted Connections
- Endpoint: GET /v1/connection/my?email=user@example.com
- Description: Get all accepted connections for a user
- Query Parameters:
email (required)

Success Response:
```json
{
  "status": "success",
  "connections": [
    {
      "id": "number",
      "name": "string",
      "email": "string",
      "role": "string"
    }
  ]
}
```

6. Get Connection Suggestions
- Endpoint: GET /v1/suggestions?email=user@example.com
- Description: Get connection suggestions based on user role and interests
- Query Parameters:
email (required)

Success Response:
```json
{
  "status": "success",
  "user_email": "string",
  "user_role": "string",
  "suggestions": [
    {
      "type": "string (investor/mentor/startup)",
      "id": "number",
      "email": "string",
      "name": "string",
      "location": "string",
      "match_reason": "string",
      "score": "number"
    }
  ]
}
```

## Startup Management

1. Get All Startups
- Endpoint: GET /v1/startup/all
- Description: Get all registered startups with entrepreneur details

Success Response:
```json
{
  "status": "success",
  "startups": [
    {
      "id": "number",
      "startup_name": "string",
      "industry": "string",
      "description": "string",
      "budget": "number",
      "timeframe": "string",
      "created_at": "timestamp",
      "updated_at": "timestamp",
      "entrepreneur": {
        "id": "number",
        "name": "string",
        "email": "string",
        "location": "string"
      }
    }
  ]
}
```

2. Get Startups by User
- Endpoint: GET /v1/startup/user?email=user@example.com
- Description: Get all startups for a specific entrepreneur
- Query Parameters:
email (required)

Success Response:

```json
{
  "status": "success",
  "startups": [
    {
      "id": "number",
      "startup_name": "string",
      "industry": "string",
      "description": "string",
      "budget": "number",
      "timeframe": "string",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
}
```

3. Insert Startup
- Endpoint: POST /v1/startup/insert
- Description: Add a new startup
- Request Body:
```json
{
  "entrepreneur_email": "string",
  "startup_name": "string",
  "industry": "string",
  "description": "string",
  "budget": "number",
  "timeframe": "string"
}
```

- Success Response:
```json
{
  "status": "success",
  "message": "Startup inserted successfully"
}
```

4. Update Startup
- Endpoint: PUT /v1/startup/update
- Description: Update an existing startup
Request Body:
```json
{
  "email": "string",
  "startup_id": "number",
  "startup_name": "string",
  "industry": "string",
  "description": "string",
  "budget": "number",
  "timeframe": "string"
}
```

- Success Response:

```json
{
  "status": "success",
  "message": "Startup updated successfully"
}
```

5. Delete Startup
- Endpoint: DELETE /v1/startup/delete
- Description: Delete a startup
- Request Body:
```json
{
  "email": "string",
  "startup_id": "number"
}
```

- Success Response:
```json
{
  "status": "success",
  "message": "Startup deleted successfully"
}
```

## Startup Engagement Features
1. Add Comment
- Endpoint: POST /v1/startup/comment
- Description: Add a comment to a startup
- Request Body:
```json
{
  "email": "string",
  "startup_id": "number",
  "content": "string"
}
```

- Success Response:
```json
{
  "status": "success",
  "message": "Comment added successfully"
}
```
2. Get Comments
- Endpoint: GET /v1/startup/comments/{id}
- Description: Get all comments for a startup
- URL Parameters:
id (startup ID)

- Success Response:
```json 
{
  "status": "success",
  "comments": [
    {
      "user_id": "number",
      "startup_id": "number",
      "content": "string"
    }
  ]
}
```
3. Like Startup
- Endpoint: POST /v1/startup/like
- Description: Like a startup
- Request Body:
```json
{
  "email": "string",
  "startup_id": "number"
}
```

- Success Response:
```json
{
  "status": "success",
  "message": "Startup liked successfully"
}
```

4. Unlike Startup
- Endpoint: POST /v1/startup/unlike
- Description: Remove like from a startup
- Request Body:
```json
{
  "email": "string",
  "startup_id": "number"
}
```

- Success Response:
```json
{
  "status": "success",
  "message": "Startup unliked successfully"
}
```

5. Get Like Count
- Endpoint: GET /v1/startup/likes/{id}
- Description: Get number of likes for a startup
- URL Parameters:
id (startup ID)

- Success Response:
```json
{
  "status": "success",
  "like_count": "number"
}
```

6. Report Startup
- Endpoint: POST /v1/startup/report
- Description: Report a startup for inappropriate content
- Request Body:
```json
{
  "email": "string",
  "startup_id": "number",
  "reason": "string"
}
```
- Success Responses:
```json
{
  "status": "success",
  "message": "Report submitted successfully"
}
```

```json
{
  "status": "success",
  "message": "Startup has been removed"
}
```

- Error Response:
```json
{
  "status": "error",
  "message": "You have already reported this startup"
}
```

## Backend Unit Tests

1. Comments Feature
- Add Comment (Success, Invalid Payload, User Not Found)
- Get Comments by Startup ID (Valid ID, Invalid ID)

2. Likes Feature
- Like Startup (Success, Already Liked, Invalid Payload)
- Unlike Startup (Success, Invalid Payload)
- Get Like Count (Correct Count Retrieval)

3. Reports Feature
- Submit Report (Success, Duplicate Report, Auto-Delete After 5 Reports, Invalid Payload, User Not Found)

4. Startup Feature
- Insert Startup (Success, User Not Found)
- Get All Startups
- Get Startups by User (Valid User, User Not Found)
- Update Startup
- Delete Startup (Success, Not Owner)

5. Suggestions Feature
- Get Suggestions Handler (Missing Email, User Not Found, Invalid Role, Valid Entrepreneur Suggestions)

### Error Responses
All endpoints return appropriate HTTP status codes with error messages in the following format:

### Technical Details
- Framework: Go with Chi router
- Database: PostgreSQL with GORM ORM
- CORS: Enabled for cross-origin requests
- Timeout: 60 seconds for requests
- Supported Methods: GET, POST, PUT, DELETE, OPTIONS

