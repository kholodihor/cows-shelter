import { useState } from 'react';
import { useAppDispatch } from '@/store/hook';
import { useTranslation } from 'react-i18next';
import { closeModal } from '@/store/slices/modalSlice';
import CloseIcon from '../icons/CloseIconMenu';
import iconCalendar from '@/assets/icons/icon_calendar.svg';
import { Post } from '@/store/slices/newsSlice';
import LoadingSpinner from '../common/LoadingSpinner';

type NewsModalProps = {
  isOpen: boolean;
  onClose: () => void;
  post: Post;
};

interface FormattedDate {
  day: string;
  month: string;
  year: string;
}

const NewsModal = ({ isOpen, onClose, post }: NewsModalProps) => {
  const { t, i18n } = useTranslation();
  const { language } = i18n;
  const dispatch = useAppDispatch();
  const [imageLoaded, setImageLoaded] = useState(false);
  const [imageError, setImageError] = useState(false);

  // Get post from props instead of Redux store

  const handleCloseModal = () => {
    onClose();
    dispatch(closeModal());
  };

  // Format the date from the post, handling both createdAt and created_at
  const formatDate = (dateString?: string): FormattedDate => {
    if (!dateString) return { day: '', month: '', year: '' };

    try {
      // Handle both string and number timestamps
      const date = new Date(dateString);

      // Check if the date is valid
      if (isNaN(date.getTime())) {
        console.error('Invalid date:', dateString);
        return { day: '', month: '', year: '' };
      }

      const day = date.getDate().toString().padStart(2, '0');
      const month = (date.getMonth() + 1).toString().padStart(2, '0');
      const year = date.getFullYear().toString();

      return { day, month, year };
    } catch (error) {
      console.error('Error formatting date:', error);
      return { day: '', month: '', year: '' };
    }
  };

  const formattedDate = formatDate(post?.createdAt);
  const displayDate =
    formattedDate.day && formattedDate.month && formattedDate.year
      ? `${formattedDate.day}.${formattedDate.month}.${formattedDate.year}`
      : '';

  const handleImageLoad = () => {
    setImageLoaded(true);
    setImageError(false);
  };

  const handleImageError = () => {
    setImageError(true);
    setImageLoaded(false);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center overflow-hidden  ">
      <div
        className={`absolute inset-0 overflow-hidden bg-black opacity-40 transition-opacity  duration-700  ${
          isOpen ? 'bg-opacity-40' : 'bg-opacity-0'
        } `}
        onClick={handleCloseModal}
      ></div>
      <div
        className={`right-50 top-50 absolute translate-x-0 transition-all duration-700 ${
          isOpen ? 'translate-x-0' : 'translate-x-full'
        } max-w-full overflow-y-auto bg-white p-6 px-5 pt-9 md:w-[672px] md:px-10 lg:w-[1136px] lg:py-16`}
        style={{
          maxHeight: isOpen ? '90vh' : 'initial'
        }}
      >
        <h2 className="text-lg font-semibold lg:divide-y-4 lg:text-2xl lg:font-bold">
          {language === 'uk' ? post.title_ua : post.title_en}
        </h2>
        <hr className=" my-2 h-px border-t-0 bg-slate-300 opacity-0 md:opacity-100" />
        <div className="py-4">
          <div className="flex items-center">
            <img
              src={iconCalendar}
              alt={t('news.calendarIconAlt') || 'Calendar icon'}
              className="mr-3"
            />
            {displayDate && (
              <span className="text-sm font-normal">
                {t('news.publishedOn', { date: displayDate })}
              </span>
            )}
          </div>
        </div>
        <div className="w-full text-justify sm:text-left lg:columns-2 lg:px-16">
          <div className="relative mb-4 h-52 w-full md:h-52 md:w-[582px] lg:h-[278px] lg:w-[488px]">
            {!imageLoaded && !imageError && (
              <div className="flex h-full items-center justify-center bg-gray-100">
                <LoadingSpinner size={24} />
              </div>
            )}
            {imageError ? (
              <div className="flex h-full items-center justify-center bg-gray-100 text-gray-500">
                {t('errors.imageLoadFailed')}
              </div>
            ) : (
              <img
                className={`mx-auto h-full w-full object-contain transition-opacity duration-300 ${
                  imageLoaded ? 'opacity-100' : 'opacity-0'
                } md:object-cover`}
                src={post.image_url}
                alt={language === 'uk' ? post.title_ua : post.title_en}
                onLoad={handleImageLoad}
                onError={handleImageError}
              />
            )}
          </div>
          <div className="prose max-w-none">
            <p className="whitespace-pre-line">
              {language === 'uk' ? post.content_ua : post.content_en}
            </p>
          </div>
        </div>
        <button
          type="button"
          onClick={handleCloseModal}
          className="absolute right-3 top-3 md:right-12 md:top-7"
        >
          <CloseIcon />
        </button>
      </div>
    </div>
  );
};

export default NewsModal;
