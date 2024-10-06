import { Box, Heading, Text, Stack, Avatar, Flex, Button, HStack } from "@chakra-ui/react";
import useUserStore from "../../stores/userStore";
import { getPosts } from "../../apiCalls/feedServiceCalls";
import { useEffect, useState } from "react";
import { Post } from "../../types/post";
import { PostCard } from "../../components/PostCard";
import { useNavigate, useParams } from "react-router-dom";
import { User } from "../../types/user";
import { followUser, getUser, unfollowUser } from "../../apiCalls/userServiceCalls";

const UserPage = () => {
  const [user, setUser] = useState<User | null>(null);

  const authUser = useUserStore((state) => state.user);
  const token = useUserStore((state) => state.token);

  const[isFollowing, setIsFollowing] = useState<boolean>(authUser != null && authUser.following.includes(user?.id!));

  const [posts, setPosts] = useState<Post[]>([]);
  const { username } = useParams<{ username: string }>(); // Get the route variable
  const navigate = useNavigate();

  const onFollowUser = async () => {
    if (!token) {
      navigate("/signin");
      return;
    }
    if (isFollowing) {
      await unfollowUser(user!.id!, token!);
      setUser({ ...user!, followers: user!.followers!.filter((id) => id !== authUser!.id) });
      setIsFollowing(false);
      return;
    } else {
      await followUser(user!.id!, token!);
      setUser({ ...user!, followers: [...user!.followers!, authUser!.id] });
      setIsFollowing(true);
    }

  }

  useEffect(() => {
    // Call the getPosts function from the feed service
    // and set the posts state
    if (authUser?.username === username) {
      navigate("/profile");
    }

    const fetchData = async () => {
      const response = await getUser(username!);
      setUser(response);
      const response2 = await getPosts(response.id!);
      setPosts(response2);
      setIsFollowing(response.followers.includes(authUser?.id!));
    };
    fetchData();
  }, [authUser, authUser?.username, navigate, username]); 

  return (
    <Box p={5} color="text" pt={20}>
      <Flex direction="column" align="center">
        <Stack spacing={3} align="center">
          <Heading as="h1" mb={5}>
            User Page
          </Heading>
          <Avatar size="xl" name={user?.username} mb={5} src={user?.avatar}/>\{" "}
          <Text fontSize="xl" fontWeight="bold">
            {user?.username}
          </Text>
          <Text fontSize="lg">
            <strong>First Name:</strong> {user?.first_name}
          </Text>
          <Text fontSize="lg">
            <strong>Last Name:</strong> {user?.last_name}
          </Text>
          <HStack spacing={5} mt={3}>
            <Text fontSize="lg">
              <strong>Following:</strong> {user?.following?.length}
            </Text>
            <Text fontSize="lg">
              <strong>Followers:</strong> {user?.followers?.length}
            </Text>
          </HStack>
          <Button mt={5} bg="primary" _hover={{ bg: "accent" }} borderRadius={20} onClick={onFollowUser}>
          {isFollowing ? "Unfollow" : "Follow"}
        </Button>
        </Stack>
        <Heading as="h1" mt={5}>
          Posts
        </Heading>
        {posts.map((post) => (
          <PostCard key={post.id} post={post} />
        ))}
      </Flex>
    </Box>
  );
};

export default UserPage;
