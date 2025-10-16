## 🛡️ User Services vs. Auth/Identity Services: A Full Markdown Guide

In a microservices architecture like the one for the RUB Bhutan colleges platform, separating **User Services** from **Auth/Identity Services** is a fundamental best practice for **security, scalability, and loose coupling**.

The core distinction is simple:

* **Auth/Identity Services** answer the question: **"Who are you, and what are you allowed to do?"** (Focus: **Security**)
* **User Services** answer the question: **"What information do we know about you?"** (Focus: **Data**)

---

### I. The Core Difference: Security vs. Data Management

| Feature | 🔑 Auth/Identity Services | 👤 User Services |
| :--- | :--- | :--- |
| **Primary Focus** | Authentication, Authorization, and Non-Repudiation. | Profile Management and User Attributes. |
| **Data Stored** | Passwords (hashed), Refresh Tokens, Access Token Keys, User Permissions/Scopes. | Name, College ID, Contact Info, Enrollment Status, User Role (as an attribute). |
| **Key Output** | A **Secure Token** (e.g., JWT) containing the User ID and permissions/roles. | The **User's Profile Data** and related attributes. |
| **Security Handling** | Highest security clearance required; handles sensitive credentials and cryptographic operations. | Standard security for PII (Personally Identifiable Information). |
| **Database** | Often a separate database optimized for quick credential lookup and token revocation. | A separate database optimized for profile data retrieval and updates. |

---

### II. Use Cases for Auth/Identity Services (Security Focus)

This service manages the *trust* relationship between the user and the system, ensuring only validated users can access protected resources.

| ID | Use Case | Actor | Description | Output/Postcondition |
| :--- | :--- | :--- | :--- | :--- |
| **AUTH.1** | **User Sign-in** | College Officer/API Gateway | Accepts user credentials, verifies the password hash against the stored record, and validates the user's active status. | Issues a signed **JWT Access Token** containing the User ID and the assigned Role (e.g., Financial Officer). |
| **AUTH.2** | **Token Validation/Introspection** | API Gateway / Any Microservice | Receives a token from an incoming request and verifies its signature, expiration, and format. | Returns `TRUE` (valid) or `FALSE` (invalid/expired); often happens on every secure API call. |
| **AUTH.3** | **Password Reset Flow** | Any User | Manages the secure process of proving identity (e.g., sending a link/code to email) to reset a forgotten password without knowing the old one. | User's new password hash is stored; old password hash is invalidated. |
| **AUTH.4** | **Role/Permission Check**** | API Gateway / Microservice | Verifies if the authenticated user (based on their token) has the necessary role or scope to execute a specific action (e.g., `can_distribute_stipend`). | Returns `Authorized` or `Unauthorized (403 Forbidden)`. |
| **AUTH.5** | **Token Revocation/Logout** | Any User | Explicitly invalidates an issued access or refresh token, terminating the active session. | The token is added to a blocklist/revocation list, preventing its future use. |

---

### III. Use Cases for User Services (Data Focus)

This service acts as the canonical source for all descriptive information about a staff member or administrator on the platform.

| ID | Use Case | Actor | Description | Output/Postcondition |
| :--- | :--- | :--- | :--- | :--- |
| **USER.1** | **Create New Staff User Record** | System Admin | Creates the initial data profile for a new college staff user (name, contact, primary college affiliation). | A new User profile record is created, including the initial role attribute (e.g., `Financial Officer`). |
| **USER.2** | **Retrieve User Profile** | Any Microservice / API Gateway | Fetches the non-sensitive profile data (name, email, role, college affiliation) based on a User ID. | Returns the full user profile data object for display or internal processing. |
| **USER.3** | **Update Profile Details** | Staff User | Allows the authenticated user to update their personal information (e.g., change of phone number or secondary email). | The specific profile fields in the database are updated. |
| **USER.4** | **Change User Role/Status** | System Admin | Modifies the role attribute of an existing user (e.g., changing a user from `Admin` to `Financial Officer`) or setting the status (e.g., `Suspended`). | The role or status attribute on the user record is updated; subsequent tokens issued by Auth Services reflect this change. |
| **USER.5** | **List All Users by Role/College** | System Admin | Provides a filtered list of all users based on administrative criteria (e.g., "List all Financial Officers at College X"). | Returns a collection of user profiles matching the criteria. |

---

### IV. The Flow Example: Accessing Financial Data

1.  **Log in:** User sends credentials to **Auth Service** (AUTH.1). Auth Service responds with a **JWT**.
2.  **Request Data:** User sends a request to the API Gateway to distribute stipends, including the **JWT**.
3.  **Authentication/Authorization:** The Gateway or **Financial Service** validates the JWT via **Auth Service** (AUTH.2). It confirms the token is valid and extracts the `User ID` and the `Role: Financial Officer` (AUTH.4).
4.  **Data Fetch:** The **Financial Service** needs the user's name for a transaction log. It calls the **User Service** using the `User ID` extracted from the token (USER.2).
5.  **Action:** The **Financial Service** proceeds with the stipend distribution business logic.

This strict separation ensures that the critical security logic (Auth Service) can be maintained and scaled independently from the profile data logic (User Service).