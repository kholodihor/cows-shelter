import React from 'react';
import { BeatLoader } from 'react-spinners';

interface LoadingSpinnerProps {
  size?: number;
  color?: string;
  className?: string;
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 15,
  color = '#4F46E5',
  className = ''
}) => {
  return (
    <div className={`flex justify-center items-center ${className}`}>
      <BeatLoader color={color} size={size} />
    </div>
  );
};

export default LoadingSpinner;
