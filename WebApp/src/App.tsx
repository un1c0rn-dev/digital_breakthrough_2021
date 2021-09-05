import React from 'react';
import { CSSTransition } from 'react-transition-group';

import { useSelector } from 'react-redux';

import { resultsSelector } from './Store/selectors';

import { Search } from './Components/Search';
import { SearchResults } from './Components/SearchResults';

import s from './App.module.scss';
import a from './Styles/Animations.module.scss';

export const App = () => {
  const results = useSelector(resultsSelector);
  return (
    <div className={s.App}>
      <Search />
      <CSSTransition in={!!results.length} timeout={200} mountOnEnter classNames={{ ...a }}>
        <SearchResults results={results} />
      </CSSTransition>
    </div>
  );
};
