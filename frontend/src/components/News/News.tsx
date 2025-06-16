import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useInView } from 'react-intersection-observer';

import { useAppDispatch, useAppSelector } from '@/store/hook';
import { openModal, closeModal } from '@/store/slices/modalSlice';
import { setActiveLink } from '@/store/slices/observationSlice';
import { fetchPosts, Post } from '@/store/slices/newsSlice';
import { useWidth } from '@/hooks/useWidth';
import LoadingSpinner from '../common/LoadingSpinner';
import ErrorMessage from '../common/ErrorMessage';
import NewsModal from '../modals/NewsModal';
import Slider from '../Slider';
import NewsBlock from './NewsBlock';

import 'swiper/css/pagination';
import 'swiper/css';

const News = () => {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();
  const screenWidth = useWidth();
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(5);
  const [pagesLength, setPagesLength] = useState(0);

  // Get news data from Redux store
  const { posts, loading, error } = useAppSelector((state) => ({
    posts: state.posts.posts,
    loading: state.posts.loading,
    error: state.posts.error
  }));

  const modalState = useAppSelector((state) => state.modals);
  const isModalOpen = modalState.isModalOpen && modalState.type === 'news';
  const selectedPost = modalState.data;
  const totalLength = posts.length;

  console.log(currentPage);

  const { ref, inView } = useInView({
    threshold: 0.5
  });

  // Fetch news data on component mount
  useEffect(() => {
    dispatch(fetchPosts());
  }, [dispatch]);

  useEffect(() => {
    if (screenWidth >= 1280) {
      setItemsPerPage(5);
    } else if (screenWidth >= 768) {
      setItemsPerPage(3);
    } else {
      setItemsPerPage(1);
    }
  }, [screenWidth, totalLength]);

  useEffect(() => {
    const pagesNumber = totalLength / itemsPerPage;
    setPagesLength(pagesNumber < 5 ? pagesNumber : 5);
  }, [totalLength, itemsPerPage]);

  // Open modal with news item
  const handleNewsItemClick = (post: Post) => {
    dispatch(openModal({ data: post, type: 'news' }));
  };

  useEffect(() => {
    if (inView) {
      dispatch(setActiveLink('#news'));
    } else {
      dispatch(setActiveLink(''));
    }
  }, [inView, dispatch]);

  // Show loading state
  if (loading && posts.length === 0) {
    return (
      <section id="news" className="bg-white py-20 md:py-40">
        <div className="container mx-auto px-4">
          <h2 className="text-main-text-color mb-10 text-2xl font-bold md:text-4xl">
            {t('news:news')}
          </h2>
          <div className="flex h-64 items-center justify-center">
            <LoadingSpinner size={30} />
          </div>
        </div>
      </section>
    );
  }

  // Show error state
  if (error) {
    return (
      <section id="news" className="bg-white py-20 md:py-40">
        <div className="container mx-auto px-4">
          <h2 className="text-main-text-color mb-10 text-2xl font-bold md:text-4xl">
            {t('news:news')}
          </h2>
          <ErrorMessage
            message={error || 'Failed to load news'}
            onRetry={() => dispatch(fetchPosts())}
          />
        </div>
      </section>
    );
  }

  // Show empty state
  if (posts.length === 0) {
    return (
      <section id="news" className="bg-white py-20 md:py-40">
        <div className="container mx-auto px-4">
          <h2 className="text-main-text-color mb-10 text-2xl font-bold md:text-4xl">
            {t('news:news')}
          </h2>
          <p className="text-center text-gray-500">
            {t('news:noNews') || 'No news available at the moment.'}
          </p>
        </div>
      </section>
    );
  }

  return (
    <section id="news" ref={ref} className="bg-white py-20 md:py-40">
      <div className="container mx-auto px-4">
        <div className="relative">
          <Slider
            title={t('news:title')}
            setCurrentPage={setCurrentPage}
            pagesLength={pagesLength}
          >
            <NewsBlock posts={posts} onItemClick={handleNewsItemClick} />
          </Slider>
        </div>

        {isModalOpen && selectedPost && (
          <NewsModal
            isOpen={isModalOpen}
            onClose={() => dispatch(closeModal())}
            post={selectedPost as Post}
          />
        )}
      </div>
    </section>
  );
};

export default News;
