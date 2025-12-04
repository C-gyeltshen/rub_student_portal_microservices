# CSV Upload Fix Guide

## Problem

You were getting the error:

```
Failed to create student data: HTTP error! status: 400, body: {"index":0,"message":"Missing required fields in student at index 0"
```

The issue was that the CSV data wasn't being correctly parsed into the JSON format expected by the backend.

## Solution

I've implemented **two solutions** for you:

### ✅ Solution 1: Direct CSV Upload (RECOMMENDED)

**New Endpoint:** `POST /api/students/bulk/csv`

This endpoint accepts CSV files directly without needing frontend parsing.

**How to use:**

1. **Using cURL:**

```bash
curl -X POST http://localhost:8086/api/students/bulk/csv \
  -F "file=@new.csv"
```

2. **Using the test script:**

```bash
./test_csv_upload.sh
```

3. **From Frontend (JavaScript/Fetch):**

```javascript
const formData = new FormData();
formData.append("file", csvFile); // csvFile is the File object from input

fetch("http://localhost:8086/api/students/bulk/csv", {
  method: "POST",
  body: formData,
})
  .then((response) => response.json())
  .then((data) => {
    console.log("Upload successful:", data);
  })
  .catch((error) => {
    console.error("Upload failed:", error);
  });
```

4. **From Frontend (HTML Form):**

```html
<form id="csvUploadForm">
  <input type="file" id="csvFile" accept=".csv" />
  <button type="submit">Upload CSV</button>
</form>

<script>
  document
    .getElementById("csvUploadForm")
    .addEventListener("submit", async (e) => {
      e.preventDefault();

      const fileInput = document.getElementById("csvFile");
      const formData = new FormData();
      formData.append("file", fileInput.files[0]);

      try {
        const response = await fetch(
          "http://localhost:8086/api/students/bulk/csv",
          {
            method: "POST",
            body: formData,
          }
        );

        const result = await response.json();

        if (result.success) {
          alert(`Successfully created ${result.created_count} students!`);
        } else {
          console.error("Validation errors:", result.validation_errors);
          alert(`Upload failed: ${result.message}`);
        }
      } catch (error) {
        console.error("Error:", error);
        alert("Upload failed: " + error.message);
      }
    });
</script>
```

### ✅ Solution 2: Improved JSON Endpoint (Existing)

**Endpoint:** `POST /api/students/bulk`

This endpoint now provides better error messages showing exactly which fields are missing.

**How to fix your frontend CSV parsing:**

If you want to keep parsing CSV on the frontend, ensure your CSV parser correctly maps field names:

```javascript
// Example CSV parsing fix
function parseCSV(csvText) {
  const lines = csvText.split("\n");
  const headers = lines[0].split(",");
  const students = [];

  for (let i = 1; i < lines.length; i++) {
    if (!lines[i].trim()) continue;

    const values = lines[i].split(",");
    const student = {};

    // Map CSV columns to JSON fields
    headers.forEach((header, index) => {
      const cleanHeader = header.trim();
      const cleanValue = values[index] ? values[index].trim() : "";

      // Handle numeric fields
      if (
        ["program_id", "college_id", "year_of_study", "semester"].includes(
          cleanHeader
        )
      ) {
        student[cleanHeader] = cleanValue ? parseInt(cleanValue) : 0;
      } else if (cleanHeader === "gpa") {
        student[cleanHeader] = cleanValue ? parseFloat(cleanValue) : 0;
      } else {
        student[cleanHeader] = cleanValue;
      }
    });

    students.push(student);
  }

  return students;
}

// Then send to the JSON endpoint
fetch("http://localhost:8086/api/students/bulk", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify(students),
});
```

## CSV File Format

Your CSV file should have these headers (exactly as shown):

```csv
student_id,first_name,last_name,email,phone_number,cid,date_of_birth,gender,permanent_address,current_address,program_id,college_id,year_of_study,semester,enrollment_date,graduation_date,gpa,status,academic_standing,guardian_name,guardian_phone_number,guardian_relation
```

**Required fields (cannot be empty):**

- `student_id`
- `first_name`
- `last_name`
- `email`

**Optional fields (can be empty):**

- All other fields

## Testing

1. **Start the service:**

```bash
cd services/student_management_service
go run main.go
```

2. **Test CSV upload:**

```bash
./test_csv_upload.sh
```

3. **Expected successful response:**

```json
{
  "success": true,
  "created_count": 5,
  "total_count": 5,
  "errors": [],
  "parse_errors": [],
  "message": "CSV bulk creation completed"
}
```

## Error Responses

### Validation Errors

```json
{
  "success": false,
  "message": "Validation failed for 1 student(s)",
  "validation_errors": [
    {
      "row": 2,
      "missing_fields": ["student_id", "first_name"],
      "received_data": {
        "student_id": "",
        "first_name": "",
        "last_name": "Lhamo",
        "email": "test@example.com"
      }
    }
  ],
  "total_rows": 5
}
```

This detailed error response will help you identify exactly which row and which fields are problematic.

## Common Issues & Solutions

### Issue: Empty required fields

**Cause:** CSV parsing shifts columns  
**Solution:** Use the direct CSV upload endpoint (`/api/students/bulk/csv`)

### Issue: BOM (Byte Order Mark) in CSV

**Cause:** CSV saved with UTF-8 BOM encoding  
**Solution:** The new CSV endpoint handles this automatically

### Issue: Column mismatch

**Cause:** Header names don't match expected field names  
**Solution:** Ensure CSV headers match exactly (case-sensitive)

## Summary

- ✅ **Use `/api/students/bulk/csv` for direct CSV upload** (recommended)
- ✅ **Better error messages** show exactly what's wrong
- ✅ **Test script** provided for easy testing
- ✅ **Frontend examples** for both approaches

The new CSV endpoint handles all the parsing server-side, eliminating frontend parsing issues!
