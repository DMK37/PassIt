import {
  Avatar,
  Box,
  HStack,
  IconButton,
  SimpleGrid,
  Text,
  VStack,
  Image,
} from "@chakra-ui/react";
import { FaComment, FaHeart, FaRetweet } from "react-icons/fa";
import { Post } from "../types/post";

interface PostCardProps {
  post: Post;
}

export const PostCard: React.FC<PostCardProps> = ({ post }) => {
  const date = new Date(post.timestamp * 1000);
  const formattedDate = date.toLocaleString();

  return (
    <Box
      key={post.id}
      p={5}
      shadow="md"
      borderWidth="1px"
      borderColor="gray"
      borderRadius="20"
      bg="white"
      width="100%"
      mt={5}
    >
      <HStack align="start">
        {/* <Avatar src={post.avatar} /> */}
        <Avatar name={post?.user.username} mb={5} src={post.user.avatar} />\{" "}
        <VStack align="start" spacing={1}>
          <HStack>
            <Text fontWeight="bold">
              {post.user.first_name} {post.user.last_name}
            </Text>
            <Text color="gray.600">@{post.user.username}</Text>
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
                  maxWidth="200px" // Set the maximum width
                  maxHeight="200px" // Set the maximum height
                />
              ))}
            </SimpleGrid>
          )}
          <HStack spacing={4} pt={2}>
            <IconButton
              aria-label="Comment"
              icon={<FaComment />}
              variant="ghost"
              size="sm"
              color="gray.500"
            />
            <IconButton
              aria-label="Retweet"
              icon={<FaRetweet />}
              variant="ghost"
              size="sm"
              color="gray.500"
            />
            <IconButton
              aria-label="Like"
              icon={<FaHeart />}
              variant="ghost"
              size="sm"
              color="gray.500"
            />
          </HStack>
        </VStack>
      </HStack>
    </Box>
  );
};
