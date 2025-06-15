/**
 * Transforms MinIO URLs from internal Docker network format to browser-accessible format
 * Replaces 'minio:9000' with 'localhost:9000' in the provided URL
 */
export const transformMinioUrl = (url: string): string => {
  if (!url) return '';

  // Check if the URL contains minio:9000 and replace it with localhost:9000
  if (url.includes('minio:9000')) {
    return url.replace('minio:9000', 'localhost:9000');
  }

  return url;
};

/**
 * Transforms all MinIO URLs in an object or array of objects
 * Recursively searches for URL strings that contain 'minio:9000' and transforms them
 */
export const transformMinioUrlsInData = <T>(data: T): T => {
  if (!data) return data;

  // Handle arrays
  if (Array.isArray(data)) {
    return data.map(item => transformMinioUrlsInData(item)) as unknown as T;
  }

  // Handle objects
  if (typeof data === 'object' && data !== null) {
    const result = { ...data } as Record<string, any>;

    for (const key in result) {
      if (Object.prototype.hasOwnProperty.call(result, key)) {
        const value = result[key];

        // If the value is a string and contains 'minio:9000', transform it
        if (typeof value === 'string' && value.includes('minio:9000')) {
          result[key] = transformMinioUrl(value);
        }
        // If the value is an object or array, recursively transform it
        else if (typeof value === 'object' && value !== null) {
          result[key] = transformMinioUrlsInData(value);
        }
      }
    }

    return result as unknown as T;
  }

  return data;
};
