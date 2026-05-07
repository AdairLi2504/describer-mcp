package description

const DescribeTool = `Takes an image from a URL, local file path, or base64 input and produces a written description of that image. 
You can optionally convert the image to base64 for inline transport, or compress it to reduce token usage — especially helpful for large images or when you're close to context limits.

- Image content description and summarization
- Analysis and interpretation of visual materials
- Visual question answering and decision support
- Content moderation and compliance checks
- Multimodal data labeling, categorization, and retrieval
- Generating alternative text for visually impaired users`

const DescribeImage = `The URL, local path, or base64 URI string of the image to be described.
- If the input is a URL, it should be a publicly accessible URL that points directly to the image file (e.g., https://example.com/image.jpg). The scheme of the URL must be included (http:// or https://).
- If the input is a local path, it should be an absolute path where the image is stored (e.g., /path/to/image.jpg).Input will be strictly converted to Base64 format before processing.
- If the input is a base64 data URI string, it should be a valid base64-encoded representation of the image data (e.g., data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD...).
When it is a local path or a base64 string, the base64 parameter will have no effect.`

const DescribeBase64 = `Whether to convert the image into a base64 data URI before processing.
This is useful when you want to pass image content inline instead of relying on direct file access or a remote URL.
When the image parameter is already a local path or a base64 string, this option has no effect.
Only when it is true or has no effect can the image be compressed.`

const DescribeCompress = `Whether to resize and compress the image before describing.
This helps reduce token usage and is especially useful for large images or when you are close to the context limit.
When the image is a local path or a base64 string, or when the base64 parameter is true, this option takes effect.
When the image is a URL and base64 is false, the MCP will return an error because the image cannot be compressed.`
