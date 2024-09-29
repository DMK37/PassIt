import React, { useEffect, useState } from "react";
import { Box, Container, Text, VStack } from "@chakra-ui/react";
import { Post } from "../../types/post";
import { getFollowingPosts } from "../../apiCalls/feedServiceCalls";
import useUserStore from "../../stores/userStore";
import { PostCard } from "../../components/PostCard";

const LandingPage: React.FC = () => {
  const [posts, setPosts] = useState<Post[]>([]);
  const token = useUserStore((state) => state.token);

  useEffect(() => {
    const fetchPosts = async () => {
      if (token) {
        const response = await getFollowingPosts(token);
        setPosts(response);
      }
    };
    fetchPosts();
  }, [token]);

  return (
    <Box>
      {/* Posts Section */}
      <Container maxW="container.md" mt={10} color="text">
        <VStack spacing={8}>
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </VStack>
        {posts.length === 0 && (
            <Text fontSize="xl" textAlign="center" mt={10}>
                No posts to show
            </Text>
        )}
      </Container>
    </Box>
  );
};

export default LandingPage;
