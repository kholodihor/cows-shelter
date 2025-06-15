import React from 'react';
import { FaExclamationCircle } from 'react-icons/fa';

interface ErrorMessageProps {
  message: string;
  className?: string;
  onRetry?: () => void;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({
  message,
  className = '',
  onRetry,
}) => {
  return (
    <div className={`flex flex-col items-center justify-center p-4 rounded-lg bg-red-50 text-red-700 ${className}`}>
      <div className="flex items-center">
        <FaExclamationCircle className="mr-2 text-xl" />
        <span className="font-medium">{message}</span>
      </div>
      {onRetry && (
        <button
          onClick={onRetry}
          className="mt-2 px-4 py-1 text-sm text-white bg-red-600 rounded hover:bg-red-700 transition-colors"
        >
          Retry
        </button>
      )}
    </div>
  );
};

export default ErrorMessage;
