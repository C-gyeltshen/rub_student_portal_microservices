# Firebase Setup Instructions

## Step 1: Create a Firebase Project

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Click "Add project"
3. Enter your project name (e.g., "rub-student-portal")
4. Enable Google Analytics (optional)
5. Click "Create project"

## Step 2: Enable Firebase Authentication

1. In your Firebase project, go to **Authentication** in the left sidebar
2. Click **Get started**
3. Go to the **Sign-in method** tab
4. Enable **Email/Password** provider:
   - Click on "Email/Password"
   - Toggle "Enable"
   - Click "Save"

## Step 3: Create a Service Account

1. Go to **Project Settings** (gear icon in left sidebar)
2. Click on the **Service accounts** tab
3. Click **Generate new private key**
4. Click **Generate key** - this downloads a JSON file
5. **Keep this file secure** - it contains sensitive credentials

## Step 4: Extract Service Account Credentials

From the downloaded JSON file, extract these values:

```json
{
  "type": "service_account",
  "project_id": "your-project-id", // ← FIREBASE_PROJECT_ID
  "private_key_id": "...",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...", // ← FIREBASE_PRIVATE_KEY
  "client_email": "firebase-adminsdk-xxxxx@your-project.iam.gserviceaccount.com", // ← FIREBASE_CLIENT_EMAIL
  "client_id": "..."
  // ... other fields
}
```

## Step 5: Set Environment Variables

Create a `.env` file in the `api-gateway` directory:

```bash
cp .env.example .env
```

Update the `.env` file with your Firebase credentials:

```env
FIREBASE_PROJECT_ID=your-actual-project-id
FIREBASE_CLIENT_EMAIL=firebase-adminsdk-xxxxx@your-actual-project.iam.gserviceaccount.com
FIREBASE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nYour actual private key here (keep the \n characters)\n-----END PRIVATE KEY-----\n"
```

**Important Notes:**

- Keep the quotes around the private key
- Preserve the `\n` characters in the private key
- Never commit the `.env` file to version control

## Step 6: Create Test Users

You can create test users in Firebase Console:

1. Go to **Authentication** → **Users** tab
2. Click **Add user**
3. Enter email and password
4. After creating the user, you'll need to set custom claims

## Step 7: Set Custom Claims for Users

You'll need to run this API Gateway with admin credentials to set custom claims for users. Create an admin script or use Firebase Functions:

```javascript
// Example: Set admin role for a user
const admin = require("firebase-admin");

// Initialize admin (this is done automatically in your API Gateway)
admin.initializeApp();

async function setUserRole(uid, role, collegeId = null) {
  const customClaims = {
    role: role,
    permissions: getDefaultPermissions(role),
  };

  if (collegeId) {
    customClaims.college_id = collegeId;
  }

  await admin.auth().setCustomUserClaims(uid, customClaims);
}

// Set admin role
await setUserRole("user-uid-here", "admin");

// Set finance officer role with college
await setUserRole("user-uid-here", "finance_officer", "college-uuid");

// Set student role with college
await setUserRole("user-uid-here", "student", "college-uuid");
```

## Step 8: Test Authentication

1. **Get a Firebase ID token**: Use Firebase Auth SDK in your frontend application to sign in and get an ID token

2. **Test with curl**:

```bash
# Test public endpoint
curl http://localhost:8080/health

# Test protected endpoint (will fail without token)
curl http://localhost:8080/profile

# Test with token
curl -H "Authorization: Bearer YOUR_FIREBASE_ID_TOKEN" http://localhost:8080/profile
```

## Step 9: Frontend Integration

In your frontend application (React, Angular, etc.), use Firebase Auth SDK:

```javascript
import { getAuth, signInWithEmailAndPassword } from "firebase/auth";

const auth = getAuth();

async function login(email, password) {
  const userCredential = await signInWithEmailAndPassword(
    auth,
    email,
    password
  );
  const idToken = await userCredential.user.getIdToken();

  // Use this token in Authorization header for API calls
  return idToken;
}
```

## Security Best Practices

1. **Never expose service account credentials** in client-side code
2. **Use HTTPS** in production
3. **Rotate service account keys** regularly
4. **Monitor authentication logs** for suspicious activity
5. **Implement proper token refresh** in your frontend
6. **Use Firebase Security Rules** for Firestore/Storage if applicable

## Available Roles

- `admin`: Full system access
- `finance_officer`: Budget and expense management, college-specific access
- `student`: Limited access, can view own data and bank details

## Testing Users

After setting up, create test users for each role:

1. **Admin User**: Can access all endpoints
2. **Finance Officer**: Can access finance endpoints and user management
3. **Student**: Can access basic endpoints and own data

## Troubleshooting

### Common Issues:

1. **"Firebase Auth client not initialized"**: Check environment variables are set correctly
2. **"Token verification failed"**: Ensure token is recent and user exists
3. **"Invalid private key"**: Check private key format and escape characters
4. **"Permission denied"**: Verify user has correct role and custom claims

### Debug Steps:

1. Check Firebase Console for user authentication logs
2. Verify environment variables are loaded: `echo $FIREBASE_PROJECT_ID`
3. Test token generation with Firebase Auth SDK
4. Check API Gateway logs for detailed error messages
