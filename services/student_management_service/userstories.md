# Student Management Service - User Stories

## Epic: Student Profile Management

### User Story 1: Register New Student

**As** an administrative officer  
**I want to** register a new student into the system  
**So that** their details are available for stipend and academic management

**Acceptance Criteria:**

- Students must be linked to an existing User account from the User_Service
- Fields include name, ID, college, program, year, contact info
- Validation ensures uniqueness of student ID and valid program code

**Priority:** High  
**Story Points:** 5

---

### User Story 2: View Student Profile

**As** a student  
**I want to** view my complete profile details  
**So that** I can verify my personal, academic, and stipend information

**Acceptance Criteria:**

- Must fetch student data via authenticated User ID
- Should show college, program, bank account, and stipend status
- Data is displayed in a clear, organized format

**Priority:** High  
**Story Points:** 3

---

### User Story 3: Update Student Profile

**As** a student  
**I want to** update my contact and bank details  
**So that** stipend payments are made to the correct account

**Acceptance Criteria:**

- Only editable fields: contact, email, bank details
- Changes are synced to Banking_Service
- Update actions are logged for audit purposes

**Priority:** High  
**Story Points:** 5

---

### User Story 4: View Program Details

**As** a student  
**I want to** view my academic program information  
**So that** I can confirm the course of study linked to my stipend

**Acceptance Criteria:**

- Must fetch data from a Program entity
- Should show program name, duration, and stipend type (if applicable)
- Display is clear and easy to understand

**Priority:** Medium  
**Story Points:** 3

---

## Epic: Stipend Management

### User Story 5: View Stipend History

**As** a student or finance officer  
**I want to** view stipend payment history  
**So that** I can track received and pending payments

**Acceptance Criteria:**

- Data fetched from Financial_Services and Banking_Services
- Display amount, date, and transaction status
- Support filtering by date range and status

**Priority:** High  
**Story Points:** 5

---

### User Story 6: Notify Stipend Credit

**As** a student or finance officer  
**I want to** receive notification when my stipend is credited  
**So that** I stay informed about transactions

**Acceptance Criteria:**

- Integration with Notification_Service (future)
- Notification triggered after successful transaction from Financial_Services
- Notification includes transaction details

**Priority:** Medium  
**Story Points:** 3

---

### User Story 7: Calculate Stipend Eligibility

**As** a financial officer  
**I want to** determine whether a student qualifies for a stipend  
**So that** only eligible students receive payments

**Acceptance Criteria:**

- Eligibility rules based on program, semester status, and academic performance
- Integration with Financial_Services for approval
- Clear indication of eligibility status and reasons

**Priority:** High  
**Story Points:** 8

---

### User Story 8: Monitor Stipend Eligibility Statistics

**As** a financial officer  
**I want to** see statistics on eligible vs ineligible students  
**So that** the university can assess budget allocation efficiency

**Acceptance Criteria:**

- Pulls data from both Student and Financial domains
- Provides visual dashboard (future scope)
- Export capabilities for reporting

**Priority:** Medium  
**Story Points:** 5

---

## Epic: Administrative Functions

### User Story 9: Deactivate or Archive Student

**As** an admin  
**I want to** deactivate students who have graduated or withdrawn  
**So that** stipend calculation excludes inactive students

**Acceptance Criteria:**

- The student's status changes to inactive
- Stipend eligibility check automatically returns false
- Deactivation date is stored for reference

**Priority:** High  
**Story Points:** 3

---

### User Story 10: Assign College

**As** an admin  
**I want to** associate each student with their college and department  
**So that** stipend can be tracked per institution

**Acceptance Criteria:**

- Each student is mapped to exactly one college entity
- Changes in college trigger re-evaluation of stipend rules
- College assignment is validated

**Priority:** High  
**Story Points:** 3

---

### User Story 11: Assign Academic Program

**As** an admin  
**I want to** assign students to an academic program (e.g., BSc IT, BEd, etc.)  
**So that** program-specific stipend policies can be applied

**Acceptance Criteria:**

- Must validate program ID and level (undergraduate, postgraduate)
- Automatically update stipend eligibility rules
- Changes are logged for audit

**Priority:** High  
**Story Points:** 5

---

### User Story 12: Generate Student Summary Report

**As** an admin  
**I want to** generate reports on active students and stipend status  
**So that** I can monitor disbursement and academic trends

**Acceptance Criteria:**

- Summaries include college-wise and program-wise data
- Exportable in CSV or PDF format
- Reports are generated within reasonable time (< 10 seconds)

**Priority:** Medium  
**Story Points:** 8

---

## Epic: System Integration

### User Story 13: Record Stipend Allocation

**As** a system  
**I want to** record the stipend amount allocated per student  
**So that** the financial team can track pending disbursements

**Acceptance Criteria:**

- Must include allocation ID, amount, date, and approval status
- Stored in the Stipend entity and linked to a student record
- Data integrity is maintained

**Priority:** High  
**Story Points:** 5

---

### User Story 14: Sync Bank Details

**As** a system  
**I want to** synchronize student bank details with Banking_Service  
**So that** stipend transfers go to verified accounts

**Acceptance Criteria:**

- Use REST or gRPC to push updates to Banking_Service
- Validate bank account number and IFSC code
- Handle sync failures gracefully

**Priority:** High  
**Story Points:** 5

---

### User Story 15: Fetch User Identity

**As** a system  
**I want to** fetch authenticated user details from User_Service  
**So that** only verified users can access student data

**Acceptance Criteria:**

- JWT-based authentication via Auth_Service
- Failsafe if the user session is invalid or expired
- Proper error handling and logging

**Priority:** High  
**Story Points:** 5

---

### User Story 16: Send Stipend Record to Financial Service

**As** a system  
**I want to** send approved stipend records to Financial_Service  
**So that** transactions can be processed for payment

**Acceptance Criteria:**

- Data includes student ID, stipend amount, and account reference
- API call confirmation is logged
- Retry mechanism for failed transfers

**Priority:** High  
**Story Points:** 5

---

---

# Epic: QA Testing for Student Management Service (SMS)

**As** a QA engineer  
**I want to** verify that the Student Management Service functions correctly, integrates with dependent services, and meets performance, security, and data integrity standards

---

## Functional Testing

### QA Story 1: Verify Student Registration

**As** a QA engineer  
**I want to** validate the student registration API  
**So that** new students can be correctly created and linked with existing users

**Acceptance Criteria:**

- Verify successful creation with valid data
- Verify rejection when the user does not exist in User_Service
- Validate input formats (e.g., CID, email)
- Check database persistence

**Test Types:** Functional, API, Validation  
**Priority:** High

---

### QA Story 2: Test Profile Viewing

**As** a QA engineer  
**I want to** verify that students can view their profile  
**So that** correct and authorized information is displayed

**Acceptance Criteria:**

- Authenticated users can access only their own data
- Unauthenticated users receive "401 Unauthorized"
- All fields match expected database values

**Test Types:** API, Security, UI validation (if applicable)  
**Priority:** High

---

### QA Story 3: Test Profile Update

**As** a QA engineer  
**I want to** test updating student contact and bank details  
**So that** data is accurately stored and synced to Banking_Service

**Acceptance Criteria:**

- Check PATCH/PUT requests update only allowed fields
- Verify integration call to Banking Service is triggered
- Audit log created for each update

**Test Types:** Integration, API, Regression  
**Priority:** High

---

### QA Story 4: Test Deactivation of Student

**As** a QA engineer  
**I want to** test the deactivation workflow  
**So that** inactive students are excluded from stipend processing

**Acceptance Criteria:**

- Check student status changes to "inactive"
- Verify stipend eligibility returns false
- Ensure Financial Service no longer processes payments

**Test Types:** Functional, Integration, Database validation  
**Priority:** High

---

### QA Story 5: Test Program Assignment

**As** a QA engineer  
**I want to** test program assignment and validation  
**So that** correct stipend policy is applied

**Acceptance Criteria:**

- Only valid programs accepted
- Program updates trigger stipend recalculation
- Audit records stored for every change

**Test Types:** API, Integration  
**Priority:** High

---

## Business Logic Testing

### QA Story 6: Validate Stipend Eligibility Logic

**As** a QA engineer  
**I want to** test stipend eligibility rules  
**So that** only qualified students receive stipends

**Acceptance Criteria:**

- Eligibility logic matches rule definitions (e.g., active enrollment, passing grades)
- Test multiple edge cases (inactive student, incomplete semester)
- Confirm response accuracy for both eligible and ineligible students

**Test Types:** Business logic validation, Unit, Integration  
**Priority:** High

---

### QA Story 7: Validate Stipend Allocation Recording

**As** a QA engineer  
**I want to** test the stipend allocation process  
**So that** allocations are correctly recorded in the system

**Acceptance Criteria:**

- Verify record creation with valid student IDs
- Confirm transaction record generated in Financial Service
- Validate rollback if allocation fails mid-process

**Test Types:** Functional, Database, Integration  
**Priority:** High

---

## Integration Testing

### QA Story 8: Test Stipend History Retrieval

**As** a QA engineer  
**I want to** verify stipend history API  
**So that** students see accurate past payment data

**Acceptance Criteria:**

- Cross-check data with Financial and Banking microservices
- Verify sorting (by date) and pagination
- Ensure performance < 2 seconds for 100+ records

**Test Types:** API, Performance  
**Priority:** Medium

---

### QA Story 9: Verify Notification Trigger

**As** a QA engineer  
**I want to** confirm notification triggers when a stipend is credited  
**So that** students receive updates about their payments

**Acceptance Criteria:**

- Notification fired after Financial Service confirms transaction success
- Test for duplicate notifications and error handling

**Test Types:** Event-driven integration, E2E  
**Priority:** Medium

---

### QA Story 10: Validate Banking Sync

**As** a QA engineer  
**I want to** verify synchronization between SMS and Banking Service  
**So that** student bank details remain accurate

**Acceptance Criteria:**

- Successful POST/PUT call updates Banking Service
- Banking Service returns valid status codes (200/400/500)
- Retry logic tested for failure cases

**Test Types:** Integration, Negative testing  
**Priority:** High

---

### QA Story 11: Verify Financial Service API Integration

**As** a QA engineer  
**I want to** confirm stipend data transfer to Financial Service  
**So that** payments are processed without loss

**Acceptance Criteria:**

- Validate payload format and required fields
- Check idempotency (no duplicate transfers)
- Ensure response handling for network delays

**Test Types:** Integration, Resilience  
**Priority:** High

---

## Security Testing

### QA Story 12: Validate Auth and User Dependency

**As** a QA engineer  
**I want to** test JWT validation and user linkage  
**So that** only authorized users can access or modify data

**Acceptance Criteria:**

- Expired tokens rejected
- Invalid or missing tokens return HTTP 401
- Correct user ownership enforced

**Test Types:** Security, API  
**Priority:** High

---

## Performance & Reporting

### QA Story 13: Verify Student Summary Report

**As** a QA engineer  
**I want to** test report generation  
**So that** admin receives accurate summaries

**Acceptance Criteria:**

- Verify data accuracy and filtering (by college, program)
- Check export formats (CSV, PDF)
- Confirm large data handling without timeout

**Test Types:** Functional, Performance, UI (if applicable)  
**Priority:** Medium

---

## Summary

**Total User Stories:** 16  
**Total QA Stories:** 13  
**Total Combined:** 29

**Priority Breakdown:**

- High Priority: 24 stories
- Medium Priority: 5 stories

**Epic Breakdown:**

- Student Profile Management: 4 stories
- Stipend Management: 4 stories
- Administrative Functions: 4 stories
- System Integration: 4 stories
- QA Testing: 13 stories
