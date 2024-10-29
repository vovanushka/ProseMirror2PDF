### Example

1. **Sample JSON Input**:
   ```json
   [
      { "type": "heading", "attrs": { "level": 1 }, "content": [{ "type": "text", "text": "Sample Document" }] },
      { "type": "paragraph", "content": [{ "type": "text", "text": "This is a paragraph with " },
                                         { "type": "text", "text": "bold", "marks": [{ "type": "bold" }] },
                                         { "type": "text", "text": " and " },
                                         { "type": "text", "text": "italic", "marks": [{ "type": "italic" }] }] },
      { "type": "bullet_list", "content": [
         { "type": "list_item", "content": [{ "type": "paragraph", "content": [{ "type": "text", "text": "First item" }] }] },
         { "type": "list_item", "content": [{ "type": "paragraph", "content": [{ "type": "text", "text": "Second item" }] }] }
      ] }
   ]


### Using Curl to Convert JSON to PDF

To send JSON data to the server and receive a PDF in response, use the following `curl` command:

```bash
curl -X POST http://localhost:8080/generate-pdf \
     -H "Content-Type: application/json" \
     -d '[{"type": "heading", "attrs": {"level": 1}, "content": [{"type": "text", "text": "Sample Document"}]},{"type": "paragraph", "content": [{"type": "text", "text": "This is a paragraph with "},{"type": "text", "text": "bold", "marks": [{"type": "bold"}]},{"type": "text", "text": " and "},{"type": "text", "text": "italic", "marks": [{"type": "italic"}]}]},{"type": "bullet_list", "content": [{"type": "list_item", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "First item"}]}]},{"type": "list_item", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Second item"}]}]}]}]' \
     --output output.pdf
