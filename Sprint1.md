## Front-end
Video Demo - https://uflorida-my.sharepoint.com/:v:/g/personal/varunaiyaswamyka_ufl_edu/Ed8bspvh5BpAssvrftRkMQ8B0AeDHBh7d1Ed_0z69Vft4A?nav=eyJyZWZlcnJhbEluZm8iOnsicmVmZXJyYWxBcHAiOiJPbmVEcml2ZUZvckJ1c2luZXNzIiwicmVmZXJyYWxBcHBQbGF0Zm9ybSI6IldlYiIsInJlZmVycmFsTW9kZSI6InZpZXciLCJyZWZlcnJhbFZpZXciOiJNeUZpbGVzTGlua0NvcHkifX0&e=gAEpXI
## User stories
- Setup Signup and Login UI
- Setup OAuth Integration

## What issues your team planned to address
- Design and create landing page with signup/signin functionalities using Next.js and React
- As a user, I want to log in using my Google or GitHub account so that I can quickly access the platform without needing to remember passwords.
- Explore OAuth options available, setup Github and Google OAuth for authentication using next-auth

## Which ones were successfully completed
- Successfully created the landing page with authentication using OAuth
- Obtained OAuth Credentials from GitHub Developer Settings and Google Cloud Console and configured .env.local with necessary credentials
- Integrated OAuth Authentication using NextAuth.js for GitHub and Google
- Configured API Route at api/auth/[...nextauth].js to handle authentication
- Implemented JWT-based Session Management with a 3-minute expiration and auto-refresh on activity
- Enabled Sign-In & Sign-Out functionality using NextAuth’s signIn and signOut methods
- Fetched User Details including email, name, and access token after authentication
- Tested Authentication Flow to ensure successful login and session handling

## Which ones didn't and why?
- Setting up LinkedIn OAuth authentication, LinkedIn OAuth requires additional permissions and approval, delaying implementation.
- Extending Session Persistence, the current session expiration is set to 3 minutes for security reasons, but adjustments may be needed for better user experience.

## Back-end
Video Demo - https://uflorida-my.sharepoint.com/:v:/g/personal/varunaiyaswamyka_ufl_edu/EfPK0_FMcDlFsMTOv77LWsMBzbITMjdufqTm0hsBWmDVxQ?e=ksyBz0&nav=eyJyZWZlcnJhbEluZm8iOnsicmVmZXJyYWxBcHAiOiJTdHJlYW1XZWJBcHAiLCJyZWZlcnJhbFZpZXciOiJTaGFyZURpYWxvZy1MaW5rIiwicmVmZXJyYWxBcHBQbGF0Zm9ybSI6IldlYiIsInJlZmVycmFsTW9kZSI6InZpZXcifX0%3D
## User stories
- Role-Based Signup:
As a user (Entrepreneur, Investor, or Mentor), I want to sign up with role-specific details so that my account reflects my needs on the platform.
- Database Management:
As a developer, I want to store user data in a structured database to ensure secure and efficient data handling.
- User Retrieval:
As an admin, I want to retrieve all registered users for platform monitoring and management.

## What issues your team planned to address
- Setup Go Backend with Gin Framework:
Initialize the Go server and configure API routes for signup and user management.
- Design and Implement Database Schema:
Create tables for Users, Entrepreneurs, Investors, and Mentors using GORM and SQLite.
- Implement Role-Based Signup API:
Develop an API endpoint to handle user registration, including role-specific data collection.
- Input Validation & Duplicate Email Handling:
Ensure that user inputs are validated, and prevent duplicate registrations using the same email.
- Develop API to Retrieve All Registered Users:
Implement an endpoint to fetch all users from the database for administrative purposes.

## Which ones were successfully completed
- Go Server Setup with Gin Framework:
Successfully initialized the server using the Gin framework.
Configured middleware for logging, error recovery, and request handling.
Database Schema Implementation:
- Designed and implemented SQLite schema using GORM.
Created tables for Users, Entrepreneurs, Investors, and Mentors.
Enabled automatic migrations for seamless schema updates.
- Role-Based Signup API:
Developed /v1/signup endpoint with role-specific data handling for Entrepreneurs, Investors, and Mentors.
Incorporated validation for user inputs and ensured correct data mapping to respective tables.
- Duplicate Email Handling:
Implemented checks to prevent users from registering with an existing email.
Returned appropriate error messages for duplicate registration attempts.
- User Retrieval API:
Developed /v1/users endpoint to retrieve all registered users from the database.
Ensured the endpoint provides a clear and structured JSON response.
- Basic Logging & Error Handling:
Integrated consistent logging for debugging and tracking server activities.
Provided clear error responses for failed API requests.

## Which ones didn't and why?
- Email Verification System:
Planned to implement an email verification process post-signup to enhance user authenticity. However, this was deprioritized due to OAuth handling verified emails and time constraints. Scheduled for consideration in future sprints if additional verification is needed.
- Role-Based Access Control (RBAC):
Planned to restrict API access based on user roles but postponed to Sprint 2. This will be crucial when implementing dashboards and role-specific functionalities.
- OAuth Token Handling on Backend:
Since OAuth login is managed by the frontend using React and NextAuth.js, backend token handling was deferred. Future integration may involve secure session validation or token-based API access.
