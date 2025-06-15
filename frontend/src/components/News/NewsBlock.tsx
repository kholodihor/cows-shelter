/* eslint-disable prettier/prettier */
import { useAppDispatch } from '@/store/hook';
import { openModal } from '@/store/slices/modalSlice';
import { Post } from '@/store/slices/newsSlice';
import { useTranslation } from 'react-i18next';
import LittleArrow from '../icons/LittleArrow';

type NewsBlockProps = {
  post: Post;
  onClick?: () => void;
};

const NewsBlock = ({ post, onClick }: NewsBlockProps) => {
  const { language } = useTranslation().i18n;
  const dispatch = useAppDispatch();

  const handleClick = () => {
    if (onClick) {
      onClick();
    } else {
      dispatch(openModal({ data: post, type: 'news' }));
    }
  };

  return (
    <div className="h-full">
      <div 
        className={`group relative cursor-pointer h-full`}
        onClick={handleClick}
      >
        <img
          src={post.image_url}
          alt={`News Image`}
          className="h-full w-full object-cover"
        />

        <div className="absolute inset-0 z-20 cursor-pointer bg-black/40 opacity-0 transition duration-300 ease-in-out md:group-hover:opacity-100"></div>
        <div className="via-opacity-30 absolute inset-0 z-30 flex cursor-pointer flex-col bg-gradient-to-b from-transparent to-black/40">
          <div className="flex h-full flex-col justify-end text-white">
            <div className="translate-y-14 space-y-3 p-4 duration-300 ease-in-out md:group-hover:translate-y-0">
              <h2 className="text-sm font-normal md:text-2xl">
                {language === 'uk' ? post.title_ua : post.title_en}
              </h2>
              <div className="text-sm opacity-0 md:group-hover:opacity-100 lg:text-sm">
                {language === 'uk' ? post.subtitle_ua : post.subtitle_en}
              </div>
            </div>
            <div className="flex items-center justify-between p-4">
              <div className="flex items-center gap-2 text-sm">
                <span className="text-white/70">
                  {new Date(post.createdAt).toLocaleDateString(language, {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                  })}
                </span>
              </div>
              <div className="group/arrow flex h-8 w-8 items-center justify-center rounded-full bg-white/10 transition-all duration-300 hover:bg-white/20">
                <LittleArrow hovered={false} />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NewsBlock;
