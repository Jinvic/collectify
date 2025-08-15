// src/utils/formatDate.js

// Simple function to format a date object or ISO string
export const formatDate = (dateString) => {
  if (!dateString) return 'N/A';
  const date = new Date(dateString);
  if (isNaN(date.getTime())) return 'Invalid Date'; // Handle invalid dates
  
  // Options for toLocaleDateString
  const options = { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    // second: '2-digit',
    // timeZoneName: 'short'
  };
  
  // Use browser's default locale
  return date.toLocaleDateString(undefined, options); 
};

// Function to format a date as YYYY-MM-DD (useful for date inputs)
export const formatDateForInput = (dateString) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return ''; // Handle invalid dates

    // Pad start is useful for ensuring two digits
    const pad = (num) => String(num).padStart(2, '0');

    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
};