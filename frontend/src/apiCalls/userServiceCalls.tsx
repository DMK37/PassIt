import axios from "axios";
import { hashPassword } from "../utils/Encryption";
import { User } from "../types/user";

const API_URL = "http://userservice-lb-2124152758.us-east-1.elb.amazonaws.com"; // Replace with your API URL
// localhost:8080
interface LoginResponse {
  token: string;
  user: User;
}

interface ResponseFetch {
    token: string;
    user: string;
}

interface LoginCredentials {
  email: string;
  password: string;
}

export const login = async (
  credentials: LoginCredentials
): Promise<LoginResponse> => {
  try {
    credentials.password = await hashPassword(credentials.password);
    const response = await axios.post<ResponseFetch>(
      `${API_URL}/login`,
      credentials
    );
    const user = JSON.parse(response.data.user);
    return { token: response.data.token, user };
  } catch (error) {
    if (axios.isAxiosError(error)) {
      // Handle Axios error
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Login failed");
    } else {
      // Handle other errors
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
};

interface SignUpResponse {
  token: string;
  user: User;
}

interface SignUpCredentials {
  username: string;
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

export const signUp = async (
  credentials: SignUpCredentials
): Promise<SignUpResponse> => {
  try {
    credentials.password = await hashPassword(credentials.password);
    const response = await axios.post<ResponseFetch>(
      `${API_URL}/users`,
      credentials
    );

    const user = JSON.parse(response.data.user);
    return { token: response.data.token, user };
  } catch (error) {
    if (axios.isAxiosError(error)) {
      // Handle Axios error
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Sign up failed");
    } else {
      // Handle other errors
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
};

export const getProfile = async (token: string): Promise<User> => {
  try {
    const response = await axios.get<User>(`${API_URL}/users`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      // Handle Axios error
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Failed to get profile");
    } else {
      // Handle other errors
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
}

  
export const getUser = async (username: string): Promise<User> => {
  try {
    const response = await axios.get<User>(`${API_URL}/users/${username}`);
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      // Handle Axios error
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Failed to get user");
    } else {
      // Handle other errors
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
}

export const followUser = async (userId: string, token: string): Promise<void> => {
  try {
    await axios.post(`${API_URL}/follow/${userId}`, null, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  } catch (error) {
    if (axios.isAxiosError(error)) {
      // Handle Axios error
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Failed to follow user");
    } else {
      // Handle other errors
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
}

export const unfollowUser = async (userId: string, token: string): Promise<void> => {
  try {
    await axios.post(`${API_URL}/unfollow/${userId}`, null, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  } catch (error) {
    if (axios.isAxiosError(error)) {
      // Handle Axios error
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Failed to follow user");
    } else {
      // Handle other errors
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
}

export const editProfile = async (
  username: string,
  firstName: string,
  lastName: string,
  token: string,
  avatar: File | null
): Promise<void> => {
  try {
    const formData = new FormData();
    formData.append("username", username);
    formData.append("first_name", firstName);
    formData.append("last_name", lastName);
    if (avatar) {
      formData.append("image", avatar);
    }

    await axios.put(
      `${API_URL}/profile/edit`,
      formData,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "multipart/form-data",
        },
      }
    );
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Failed to edit profile");
    } else {
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
}