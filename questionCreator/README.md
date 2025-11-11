# Question Creator Service

An AI-powered service that generates quiz questions based on topics using Gemini AI.

## API Endpoints

### POST /questions

Generates quiz questions on a specified topic using AI.

**Content-Type**: `application/json`

#### Request Body Format

```json
{
  "questionText": "kubernetes etcd",
  "numincorrectans": 3,
  "numcorrectans": 1,
  "numquestions": 5,
  "category": "kubernetes",
  "subcategory": "storage"
}
```

#### Field Descriptions

- **`questionText`** (string, required): The topic or question you want AI to generate questions about
- **`numincorrectans`** (int): Number of incorrect answers per question
- **`numcorrectans`** (int): Number of correct answers per question
- **`numquestions`** (int): How many questions to generate on this topic
- **`category`** (string, optional): Main category (if not provided, will be extracted from questionText)
- **`subcategory`** (string, optional): Sub-category (if not provided, will be extracted from questionText)

#### Example Request

```bash
curl -X POST http://localhost:8080/questions \
  -H "Content-Type: application/json" \
  -d '{
    "questionText": "Docker containers and orchestration",
    "numincorrectans": 3,
    "numcorrectans": 1,
    "numquestions": 5,
    "category": "DevOps",
    "subcategory": "Containerization"
  }'
```

#### Example Response

```json
{
  "questions": [
    {
      "questionText": "What is the primary purpose of Docker?",
      "type": "multiple-choice",
      "category": "DevOps",
      "subcategory": "Containerization",
      "dateadded": "2024-11-11",
      "certification": "",
      "answerReference": "",
      "answers": [
        {"text": "To containerize applications", "iscorrect": true},
        {"text": "To manage databases", "iscorrect": false},
        {"text": "To create web servers", "iscorrect": false},
        {"text": "To compile code", "iscorrect": false}
      ]
    }
  ]
}
```

### GET /health

Health check endpoint.

**Response**:
```json
{
  "status": "OK"
}
```

### GET /questions/health

Questions service health check endpoint.

## How It Works

1. **AI Generation**: The service uses Gemini AI to generate questions based on the provided topic
2. **File Output**: Generated questions are automatically saved to JSON files in the configured questions directory
3. **API Response**: Returns the generated questions in a structured format ready for use in quiz applications

## Configuration

The service uses configuration from environment variables or config files. Make sure to set up your Gemini AI credentials and questions output directory.

## File Structure

Generated questions are saved as JSON files in the format:
- Filename: `{category}-{subcategory}-{uuid}.json`
- Content: Array of question objects with answers

## Integration

This service is designed to work with the main Quiz Application. Generated questions can be:
- Loaded directly by the quiz app via URL
- Stored in question packs for offline use
- Used to populate the quiz database
