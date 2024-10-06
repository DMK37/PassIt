import {
  Avatar,
  Box,
  HStack,
  IconButton,
  SimpleGrid,
  Text,
  VStack,
  Image,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalCloseButton,
  ModalBody,
} from "@chakra-ui/react";
import { FaComment, FaHeart } from "react-icons/fa";
import { Post } from "../types/post";
import { useEffect, useState } from "react";
import { likePost, unlikePost } from "../apiCalls/feedServiceCalls";
import useUserStore from "../stores/userStore";
import { useNavigate } from "react-router-dom";

interface PostCardProps {
  post: Post;
}

export const PostCard: React.FC<PostCardProps> = ({ post }) => {
  const date = new Date(post.timestamp * 1000);
  const formattedDate = date.toLocaleString();
  const [isHeartClicked, setIsHeartClicked] = useState(false);
  const token = useUserStore((state) => state.token);
  const navigate = useNavigate();
  const user = useUserStore((state) => state.user);
  const [likeCount, setLikeCount] = useState(post.likes.length || 0);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedImage, setSelectedImage] = useState("");

  useEffect(() => {
    if (user?.id && post.likes.includes(user.id)) {
      setIsHeartClicked(true);
    }
  }, [post.likes, user?.id]);

  const handleHeartClick = async () => {
    if (!token) {
      navigate("/signin");
      return;
    }
    try {
      if (isHeartClicked) {
        await unlikePost(post.id, post.user_id, token);
        setLikeCount(likeCount - 1);
      } else {
        await likePost(post.id, post.user_id, token);
        setLikeCount(likeCount + 1);
      }
      setIsHeartClicked(!isHeartClicked);
    } catch (error) {
      console.error("Failed to like post:", error);
    }
  };

  const handleImageClick = (image: string) => {
    setSelectedImage(image);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setSelectedImage("");
  };

  const handleRedirectToProfile = () => {
    navigate(`/user/${post.user.username}`);
  };

  return (
    <Box
      key={post.id}
      p={5}
      borderWidth="1px"
      borderColor="gray"
      borderRadius="20"
      bg="white"
      width="100%"
      mt={5}
    >
      <HStack align="start">
        <Avatar
          name={post?.user.username}
          mb={5}
          src={post.user.avatar}
          cursor="pointer"
          onClick={handleRedirectToProfile}
        />
        \{" "}
        <VStack align="start" spacing={1}>
          <HStack>
            <Text
              fontWeight="bold"
              cursor="pointer"
              onClick={handleRedirectToProfile}
            >
              {post.user.first_name} {post.user.last_name}
            </Text>
            <Text
              color="gray.600"
              cursor="pointer"
              onClick={handleRedirectToProfile}
            >
              @{post.user.username}
            </Text>
            <Text color="gray.500">· {formattedDate}</Text>
          </HStack>

          <Text>{post.text}</Text>
          {post.images.length > 0 && (
            <SimpleGrid columns={[1, 2, 3]} spacing={2} mt={2}>
              {post.images.map((image, index) => (
                <Image
                  key={index}
                  src={image}
                  alt={`Post image ${index + 1}`}
                  borderRadius="md"
                  maxWidth="200px"
                  maxHeight="200px"
                  onClick={() => handleImageClick(image)}
                  cursor="pointer"
                />
              ))}
            </SimpleGrid>
          )}
          <HStack spacing={4} pt={2}>
            <HStack>
              <IconButton
                aria-label="Like"
                icon={<FaHeart />}
                variant="ghost"
                size="xl"
                color={isHeartClicked ? "red.500" : "gray.500"}
                onClick={handleHeartClick}
              />
              <Text>{likeCount}</Text>
            </HStack>
            <IconButton
              aria-label="Comment"
              icon={<FaComment />}
              variant="ghost"
              size="xl"
              color="gray.500"
            />
          </HStack>
        </VStack>
      </HStack>

      <Modal isOpen={isModalOpen} onClose={closeModal} size="xl">
        <ModalOverlay />
        <ModalContent>
          <ModalCloseButton />
          <ModalBody>
            <Image
              src={selectedImage}
              alt="Selected post image"
              borderRadius="md"
            />
          </ModalBody>
        </ModalContent>
      </Modal>
    </Box>
  );
};
