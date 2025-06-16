import { configureStore, Reducer, AnyAction } from '@reduxjs/toolkit';
import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';

// Import all reducers except the one causing circular dependency
import newsReducer from './slices/newsSlice';
import modalReducer from './slices/modalSlice';
import galleryReducer from './slices/gallerySlice';
import pdfReducer from './slices/pdfSlice';
import observationReducer from './slices/observationSlice';
import excursionsReducer from './slices/excursionsSlice';
import reviewsReducer from './slices/reviewsSlice';
import partnersReducer from './slices/partnersSlice';
import responseAlertReducer from './slices/responseAlertSlice';

// Debug logger middleware
const logger = (store: any) => (next: any) => (action: any) => {
  console.group(action.type);
  console.info('dispatching', action);
  const result = next(action);
  console.log('next state', store.getState());
  console.groupEnd();
  return result;
};

// Import auth reducer
import authReducer from './slices/authSlice';

// Define the state shape
type AppState = {
  posts: ReturnType<typeof newsReducer>;
  modals: ReturnType<typeof modalReducer>;
  gallery: ReturnType<typeof galleryReducer>;
  pdf: ReturnType<typeof pdfReducer>;
  observer: ReturnType<typeof observationReducer>;
  excursions: ReturnType<typeof excursionsReducer>;
  reviews: ReturnType<typeof reviewsReducer>;
  partners: ReturnType<typeof partnersReducer>;
  contacts: any; // This will be typed after store creation
  alert: ReturnType<typeof responseAlertReducer>;
  auth: ReturnType<typeof authReducer>;
};

// Create a root reducer with lazy loading for the contacts reducer
const createRootReducer = (): Reducer<AppState, AnyAction> => {
  let contactsReducer: Reducer<AppState['contacts'], AnyAction> | null = null;

  // This will be called the first time the contacts reducer is needed
  const getContactsReducer = (): Reducer<AppState['contacts'], AnyAction> => {
    if (!contactsReducer) {
      // Use dynamic import to avoid circular dependency
      import('./slices/contactsSlice').then((module) => {
        contactsReducer = module.default;
      });
      // Return a no-op reducer until the real one is loaded
      return (state = {}, _action) => state as any;
    }
    return contactsReducer!;
  };

  // Return the root reducer function
  return (state, _action) => {
    if (!state) {
      return {
        posts: newsReducer(undefined, _action),
        modals: modalReducer(undefined, _action),
        gallery: galleryReducer(undefined, _action),
        pdf: pdfReducer(undefined, _action),
        observer: observationReducer(undefined, _action),
        excursions: excursionsReducer(undefined, _action),
        reviews: reviewsReducer(undefined, _action),
        partners: partnersReducer(undefined, _action),
        contacts: getContactsReducer()(undefined, _action),
        alert: responseAlertReducer(undefined, _action),
        auth: authReducer(undefined, _action)
      } as AppState;
    }

    // Handle the action with all reducers
    return {
      posts: newsReducer(state.posts, _action),
      modals: modalReducer(state.modals, _action),
      gallery: galleryReducer(state.gallery, _action),
      pdf: pdfReducer(state.pdf, _action),
      observer: observationReducer(state.observer, _action),
      excursions: excursionsReducer(state.excursions, _action),
      reviews: reviewsReducer(state.reviews, _action),
      partners: partnersReducer(state.partners, _action),
      contacts: getContactsReducer()(state.contacts, _action),
      alert: responseAlertReducer(state.alert, _action),
      auth: authReducer(state.auth, _action)
    } as AppState;
  };
};

// Create the store with the root reducer
const store = configureStore({
  reducer: createRootReducer(),
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false
    }).concat(logger),
  devTools: process.env.NODE_ENV !== 'production'
});

// Export the store
export { store };

// Export types
export type RootState = AppState;
export type AppDispatch = typeof store.dispatch;

// Export hooks with proper typing
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
