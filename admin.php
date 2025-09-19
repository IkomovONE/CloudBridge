<?php
session_start();

// Simple hardcoded credentials
$admin_username = "admin";
$admin_password = "password123"; // Replace with your desired password

if ($_SERVER["REQUEST_METHOD"] === "POST") {
    $username = $_POST["username"] ?? '';
    $password = $_POST["password"] ?? '';

    if ($username === $admin_username && $password === $admin_password) {
        $_SESSION["admin_logged_in"] = true;
        header("Location: admin_dashboard.php"); // Redirect to admin dashboard
        exit;
    } else {
        $error = "Invalid credentials!";
    }
}
?>

<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Admin Login</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background: #f4f4f4;
      display: flex;
      justify-content: center;
      align-items: center;
      height: 100vh;
    }
    .login-box {
      background: white;
      padding: 30px;
      border-radius: 8px;
      box-shadow: 0 0 10px rgba(0,0,0,0.1);
      width: 300px;
    }
    h2 {
      text-align: center;
      color: #333;
    }
    input[type="text"], input[type="password"] {
      width: 100%;
      padding: 10px;
      margin: 10px 0 20px;
      border: 1px solid #ccc;
      border-radius: 5px;
    }
    button {
      width: 100%;
      padding: 10px;
      background-color: #008080;
      color: white;
      border: none;
      border-radius: 5px;
    }
    .error {
      color: red;
      text-align: center;
      margin-bottom: 10px;
    }
  </style>
</head>
<body>

<div class="login-box">
  <h2>Admin Login</h2>
  <?php if (!empty($error)): ?>
    <div class="error"><?= $error ?></div>
  <?php endif; ?>
  <form method="POST">
    <input type="text" name="username" placeholder="Username" required>
    <input type="password" name="password" placeholder="Password" required>
    <button type="submit">Login</button>
  </form>
</div>

</body>
</html>
