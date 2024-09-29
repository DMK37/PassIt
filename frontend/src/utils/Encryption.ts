import { hash, compare } from 'bcrypt-ts';

const FIXED_SALT = "$2a$10$fixedsaltfixedsaltfixedsa"; // You can adjust the salt rounds as needed

export const hashPassword = async (password: string): Promise<string> => {
  try {
    const hashedPassword = await hash(password, FIXED_SALT);
    return hashedPassword;
  } catch (error) {
    console.error('Error hashing password:', error);
    throw new Error('Failed to hash password');
  }
};

export const comparePassword = async (password: string, hashedPassword: string): Promise<boolean> => {
  try {
    const isMatch = await compare(password, hashedPassword);
    return isMatch;
  } catch (error) {
    console.error('Error comparing password:', error);
    throw new Error('Failed to compare password');
  }
};