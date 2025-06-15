import CloseIcon from '@/components/icons/CloseIconMenu';
import { ForgotPasswordProps } from '@/types';
import { Link } from 'react-router-dom';

const PopUpForgotPassword = ({ token, closePopup }: ForgotPasswordProps) => {
  return (
    <div
      onClick={closePopup}
      className="fixed left-0 top-0 z-50 flex h-screen w-screen items-center justify-center bg-black bg-opacity-40"
    >
      <div className="relative h-[180px] w-[492px]  bg-white px-[51px] py-[56px] ">
        <button
          className="absolute right-5 top-5 text-black"
          onClick={closePopup}
        >
          <CloseIcon />
        </button>
        <p className=" text-[1.063rem] font-medium">
          Забули пароль? Перейдіть за лінком{' '}
          <Link to={`/reset/${token}`} className="text-darkyellow">
            лінк
          </Link>{' '}
          для зміни паролю.
        </p>
      </div>
    </div>
  );
};
export default PopUpForgotPassword;
