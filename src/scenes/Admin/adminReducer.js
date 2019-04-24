import { reducer as formReducer } from 'redux-form';
import { applyMiddleware, combineReducers, compose, createStore } from 'redux';
import { routerMiddleware, routerReducer } from 'react-router-redux';
import { adminReducer, defaultI18nProvider, i18nReducer, formMiddleware, USER_LOGOUT } from 'react-admin';

export default ({ dataProvider, i18nProvider = defaultI18nProvider, history, locale = 'en' }) => {
  const reducer = combineReducers({
    admin: adminReducer,
    i18n: i18nReducer(locale, i18nProvider(locale)),
    form: formReducer,
    router: routerReducer,
  });

  const resettableAppReducer = (state, action) => reducer(action.type !== USER_LOGOUT ? state : undefined, action);

  const store = createStore(
    resettableAppReducer,
    {
      /* set your initial state here */
    },
    compose(
      applyMiddleware(
        formMiddleware,
        routerMiddleware(history),
        // add your own middlewares here
      ),
      typeof window !== 'undefined' && window.__REDUX_DEVTOOLS_EXTENSION__
        ? window.__REDUX_DEVTOOLS_EXTENSION__()
        : f => f,
      // add your own enhancers here
    ),
  );
  return store;
};
