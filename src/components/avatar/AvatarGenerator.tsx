import { Avatar } from 'antd';

interface AvatarGeneratorProps {
  userID: string;
  size?: 'small' | 'default' | 'large' | number;
}

export const AvatarGenerator: React.FC<AvatarGeneratorProps> = ({ 
  userID, 
  size = 'large' 
}) => {
  const seed = userID || 'default';
  const avatarUrl = `https://api.dicebear.com/7.x/miniavs/svg?seed=${seed}`;

  return <Avatar src={avatarUrl} size={size} />;
};
