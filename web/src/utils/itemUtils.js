// src/utils/itemUtils.js

// Map field type numbers to human-readable names
export const getFieldTypeName = (type) => {
  switch (type) {
    case 1: return 'String';
    case 2: return 'Integer';
    case 3: return 'Boolean';
    case 4: return 'Datetime';
    default: return 'Unknown';
  }
};

// Format field value based on its type
export const formatFieldValue = (value, fieldType) => {
  if (value === null || value === undefined) {
    return 'N/A';
  }

  switch (fieldType) {
    case 1: // String
      return String(value);
    case 2: // Integer
      return Number.isInteger(value) ? value.toString() : 'Invalid Integer';
    case 3: // Boolean
      return value ? 'Yes' : 'No';
    case 4: // Datetime
      // Assuming value is an ISO string
      if (typeof value === 'string') {
        const date = new Date(value);
        if (isNaN(date.getTime())) {
          return 'Invalid Date';
        }
        // You can use the formatDate utility here if you import it
        // For simplicity, using toLocaleString
        return date.toLocaleString();
      }
      return 'Invalid Date Type';
    default:
      return String(value); // Fallback
  }
};