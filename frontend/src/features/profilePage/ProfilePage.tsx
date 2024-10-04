import {
  Box,
  Heading,
  Text,
  Stack,
  Avatar,
  Button,
  Flex,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  Textarea,
  Input,
  SimpleGrid,
  ModalFooter,
  Image,
  useDisclosure,
  HStack,
} from "@chakra-ui/react";
import useUserStore from "../../stores/userStore";
import { createPost, getPosts } from "../../apiCalls/feedServiceCalls";
import { useEffect, useState } from "react";
import { Post } from "../../types/post";
import { PostCard } from "../../components/PostCard";
import { editProfile, getProfile } from "../../apiCalls/userServiceCalls";

const ProfilePage = () => {
  const user = useUserStore((state) => state.user);
  const token = useUserStore((state) => state.token);
  const setUser = useUserStore((state) => state.setUser);

  const [posts, setPosts] = useState<Post[]>([]);
  const [postText, setPostText] = useState<string>("");
  const [postImages, setPostImages] = useState<File[]>([]);
  const [username, setUsername] = useState(user?.username || "");
  const [firstName, setFirstName] = useState(user?.first_name || "");
  const [lastName, setLastName] = useState(user?.last_name || "");
  const [avatar, setAvatar] = useState<File | null>(null);

  const {
    isOpen: isPostModalOpen,
    onOpen: onPostModalOpen,
    onClose: onPostModalClose,
  } = useDisclosure();
  const {
    isOpen: isEditModalOpen,
    onOpen: onEditModalOpen,
    onClose: onEditModalClose,
  } = useDisclosure();

  useEffect(() => {
    // Call the getPosts function from the feed service
    // and set the posts state
    const fetchPosts = async () => {
      const response = await getPosts(user?.id!);
      setPosts(response);
      const response2 = await getProfile(token!);
      setUser(response2, token!);
    };
    fetchPosts();
  }, [setUser, token, user?.id, user?.username]);

  const handleSaveProfile = async () => {
    await editProfile(username, firstName, lastName, token!, avatar);
    onEditModalClose();
  };

  const handleCreatePost = async () => {
    try {
      const newPost = await createPost(postText, postImages, token!);
      setPosts([newPost, ...posts]);
      onPostModalClose();
    } catch (error) {
      console.error("Failed to create post:", error);
    }
  };

  const handleImageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files || []);
    setPostImages(files);
  };

  return (
    <Box p={5} color="text" pt={20}>
      <Flex direction="column" align="center">
        <Stack spacing={3} align="center">
          <Heading as="h1" mb={5}>
            Profile Page
          </Heading>
          <Avatar size="xl" name={user?.username} mb={5} src={user?.avatar} />\{" "}
          <Text fontSize="xl" fontWeight="bold">
            {user?.username}
          </Text>
          {/* <Text fontSize="lg">
            <strong>Email:</strong> {user?.email}
          </Text> */}
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
        </Stack>
        <Button
          mt={5}
          bg="primary"
          _hover={{ bg: "accent" }}
          borderRadius={20}
          onClick={onEditModalOpen}
          color={"background"}
        >
          Edit Profile
        </Button>
        <Heading as="h1" mt={5}>
          Posts
        </Heading>
        <Button
          mt={5}
          bg="primary"
          _hover={{ bg: "accent" }}
          borderRadius={20}
          onClick={onPostModalOpen}
          color={"background"}
        >
          Add Post
        </Button>
        {posts.map((post) => (
          <PostCard key={post.id} post={post} />
        ))}
      </Flex>
      <Modal isOpen={isPostModalOpen} onClose={onPostModalClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Create a New Post</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Textarea
              placeholder="Write your post..."
              value={postText}
              onChange={(e) => setPostText(e.target.value)}
              mb={3}
            />
            <Input
              type="file"
              accept="image/*"
              multiple
              onChange={handleImageChange}
            />
            {postImages.length > 0 && (
              <SimpleGrid columns={[1, 2, 3]} spacing={2} mt={2}>
                {postImages.map((image, index) => (
                  <Image
                    key={index}
                    src={URL.createObjectURL(image)}
                    alt={`Post image ${index + 1}`}
                    borderRadius="md"
                  />
                ))}
              </SimpleGrid>
            )}
          </ModalBody>
          <ModalFooter>
            <Button colorScheme="blue" mr={3} onClick={handleCreatePost}>
              Post
            </Button>
            <Button variant="ghost" onClick={onPostModalClose}>
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
      <Modal isOpen={isEditModalOpen} onClose={onEditModalClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Edit Profile</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Stack spacing={3}>
              <Input
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
              <Input
                placeholder="First Name"
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
              />
              <Input
                placeholder="Last Name"
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
              />
              <Input
                type="file"
                accept="image/*"
                onChange={(e) => {
                  if (e.target.files) {
                    setAvatar(e.target.files[0]);
                  }
                }}
              />
            </Stack>
          </ModalBody>
          <ModalFooter>
            <Button colorScheme="blue" onClick={handleSaveProfile} mr={3}>
              Save
            </Button>
            <Button variant="ghost" onClick={onEditModalClose}>
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Box>
  );
};

export default ProfilePage;
