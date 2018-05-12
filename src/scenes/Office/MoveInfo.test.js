import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import MoveInfo from './MoveInfo';
import store from 'shared/store';

const dummyFunc = () => {};
const moveIsLoading = false;
const moveHasLoadError = false;
const moveHasLoadSuccess = null;
const match = {
  params: { moveID: '123456' },
  url: 'www.nino.com',
  path: '/moveIt/moveIt',
};

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <Provider store={store}>
      <MoveInfo
        moveIsLoading={moveIsLoading}
        moveHasLoadError={moveHasLoadError}
        moveHasLoadSuccess={moveHasLoadSuccess}
        match={match}
        loadMove={dummyFunc}
      />
    </Provider>,
    div,
  );
});

// import React from 'react';
// import { shallow } from 'enzyme';
// import { MoveInfo } from '.';
// import Header from 'shared/Header/Office';

// describe('MoveInfo tests', () => {
//   let _moveInfo;

//   beforeEach(() => {
//     _moveInfo = shallow(<MoveInfo />);
//   });

//   it('renders without crashing', () => {
//     const moveInfo = _moveInfo.find('div');
//     expect(moveInfo).toBeDefined;
//   });

//   it('renders Header component', () => {
//     expect(_moveInfo.find(Header)).toHaveLength(1);
//   });

// });
