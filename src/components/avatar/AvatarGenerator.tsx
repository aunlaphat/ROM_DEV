import { Avatar } from 'antd';

interface AvatarGeneratorProps {
  userID: string;
  userName: string;
  size?: 'small' | 'default' | 'large' | number;
}

// ฟังก์ชันสำหรับสร้างสีที่คงที่สำหรับแต่ละ user
const generateColor = (userName: string): string => {
  // ใช้ userName เป็น seed ในการสร้างสี
  const hash = userName.split('').reduce((acc, char) => {
    return char.charCodeAt(0) + ((acc << 5) - acc);
  }, 0);
  
  // สร้างสีในรูปแบบ HSL เพื่อให้ได้สีที่สวยงาม
  const h = Math.abs(hash % 360);  // สี (0-360)
  const s = 70;  // ความอิ่มตัว (%)
  const l = 65;  // ความสว่าง (%)
  
  return `hsl(${h}, ${s}%, ${l}%)`;
};

export const AvatarGenerator: React.FC<AvatarGeneratorProps> = ({ 
  userName,
  size = 'large'
}) => {
  const initial = userName.charAt(0).toUpperCase();
  const backgroundColor = generateColor(userName);
  
  return (
    <Avatar 
      size={size} 
      style={{ 
        backgroundColor,
        color: '#ffffff',  // สีตัวอักษร
        fontWeight: 'bold'
      }}
    >
      {initial}
    </Avatar>
  );
};