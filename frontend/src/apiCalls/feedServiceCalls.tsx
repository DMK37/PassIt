import axios from "axios";
import { Post } from "../types/post";

// const API_URL = "http://userservice-lb-2124152758.us-east-1.elb.amazonaws.com";

const API_URL = "http://localhost:8081";
export const getPosts = async (userId: string): Promise<Post[]> => {
  try {
    const response = await axios.get<Post[]>(`${API_URL}/posts/users/${userId}`);
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Failed to get posts");
    } else {
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
};

export const createPost = async (
  text: string,
  images: File[],
  token: string
): Promise<Post> => {
  try {
    const formData = new FormData();
    formData.append("text", text);
    images.forEach((image) => {
      formData.append("images", image);
    });
    const response = await axios.post<Post>(`${API_URL}/posts`, formData, {
      headers: {
        "Content-Type": "multipart/form-data",
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Failed to create post");
    } else {
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
};

export const getPost = async (
  userId: string,
  postId: string
): Promise<Post> => {
  try {
    const response = await axios.get<Post>(
      `${API_URL}/posts/${postId}/users/${userId}`
    );
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Failed to get post");
    } else {
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
};

export const getFollowingPosts = async (token: string): Promise<Post[]> => {
  try {
    const response = await axios.get<Post[]>(`${API_URL}/posts-following`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error("Axios error:", error.response?.data);
      throw new Error(
        error.response?.data?.message || "Failed to get following posts"
      );
    } else {
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
}

export const likePost = async (postId: string, ownerId: string, token: string): Promise<void> => {
  try {
    await axios.post(`${API_URL}/posts/like`, 
      {
        post_id: postId,
        owner_id: ownerId,
      }, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Failed to like post");
    } else {
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
}

export const unlikePost = async (postId: string, ownerId: string, token: string): Promise<void> => {
  try {
    await axios.post(`${API_URL}/posts/unlike`, 
      {
        post_id: postId,
        owner_id: ownerId,
      }, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error("Axios error:", error.response?.data);
      throw new Error(error.response?.data?.message || "Failed to like post");
    } else {
      console.error("Unexpected error:", error);
      throw new Error("An unexpected error occurred");
    }
  }
}
