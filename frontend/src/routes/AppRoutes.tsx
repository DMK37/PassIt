import { Route, Routes } from "react-router-dom";
import LandingPage from "../features/landingPage/LandingPage";
import SignInPage from "../features/signInPage/SignInPage";
import Layout from "../components/Layout";
import SignUpPage from "../features/signUpPage/SignUpPage";
import ProfilePage from "../features/profilePage/ProfilePage";
import UserPage from "../features/userPage/userPage";

const AppRoutes: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<LandingPage />} />
        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/user/:username" element={<UserPage />} />
      </Route>
      <Route path="/signin" element={<SignInPage />} />
      <Route path="/signup" element={<SignUpPage />} />
    </Routes>
  );
};

export default AppRoutes;
