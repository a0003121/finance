function logout() {
    sessionStorage.removeItem('authToken');
    window.location.href = 'login.html';
  }