import { useState } from 'react';
import { useForm, SubmitHandler } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { singInSchema } from './schema/sindInSchema';
import { FormValuesSignIn } from '@/types';
import { useAuth } from '@/contexts/AuthContext';

interface ApiError extends Error {
  response?: {
    status: number;
    data?: {
      message?: string;
    };
  };
}

const SingInForm = () => {
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isValid }
  } = useForm<FormValuesSignIn>({
    resolver: zodResolver(singInSchema),
    mode: 'onChange'
  });

  const onSubmit: SubmitHandler<FormValuesSignIn> = async (data) => {
    setIsLoading(true);
    try {
      await login({
        email: data.email,
        password: data.password
      });
      // Navigation happens automatically after successful login in the AuthProvider
    } catch (error) {
      const apiError = error as ApiError;
      const errorMessage = apiError?.response?.data?.message || 'An error occurred during login';
      
      if (apiError.response?.status === 404) {
        setError('email', {
          type: 'manual',
          message: errorMessage
        });
      } else {
        setError('password', {
          type: 'manual',
          message: errorMessage
        });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleForgotPassword = () => {
    // TODO: Implement forgot password functionality
    console.log('Forgot password clicked');
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="mx-auto mb-[60px] w-[330px]"
    >
      <div className="mb-4">
        <label htmlFor="email" className="block mb-2">
          Email:
        </label>
        <input
          id="email"
          type="email"
          {...register('email')}
          className={`w-full border ${
            errors.email ? 'border-red' : 'border-darkgray'
          } px-4 py-2 rounded`}
          placeholder="Введіть свій email"
        />
        {errors.email && (
          <p className="text-red-500 text-sm mt-1">{errors.email.message}</p>
        )}
      </div>

      <div className="mb-6">
        <label htmlFor="password" className="block mb-2">
          Пароль:
        </label>
        <input
          id="password"
          type="password"
          {...register('password')}
          className={`w-full border ${
            errors.password ? 'border-red' : 'border-darkgray'
          } px-4 py-2 rounded`}
          placeholder="Введіть пароль"
        />
        {errors.password && (
          <p className="text-red-500 text-sm mt-1">{errors.password.message}</p>
        )}
      </div>

      <button
        type="button"
        onClick={handleForgotPassword}
        className="text-blue-500 hover:underline mb-6 block"
      >
        Не пам'ятаю пароль
      </button>

      <button
        type="submit"
        disabled={!isValid || isLoading}
        className={`w-full bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition-colors ${
          (!isValid || isLoading) ? 'opacity-50 cursor-not-allowed' : ''
        }`}
      >
        {isLoading ? 'Вхід...' : 'Увійти'}
      </button>
    </form>
  );
};

export default SingInForm;
