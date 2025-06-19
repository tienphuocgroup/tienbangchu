// Vietnamese Number Converter Frontend
// Handles user input and calls the backend API to get the Vietnamese representation

const numberInput = document.getElementById('numberInput');
const resultDiv = document.getElementById('result');

// Simple debounce utility to avoid spamming the API
function debounce(fn, delay = 300) {
  let timeout;
  return (...args) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => fn(...args), delay);
  };
}

async function fetchConversion(numberStr) {
  if (!numberStr) {
    resultDiv.textContent = '';
    return;
  }

  try {
    const res = await fetch(`/api/v1/convert?number=${numberStr}`);
    const data = await res.json();

    if (!res.ok) {
      resultDiv.innerHTML = `<span class="error">${data.error || 'Lỗi'}</span>`;
      return;
    }

    resultDiv.textContent = data.vietnamese;
  } catch (err) {
    resultDiv.innerHTML = '<span class="error">Network error</span>';
  }
}

numberInput.addEventListener(
  'input',
  debounce((e) => {
    // Keep only digits
    const sanitized = e.target.value.replace(/[^0-9]/g, '');
    e.target.value = sanitized;
    fetchConversion(sanitized);
  }, 250)
);
