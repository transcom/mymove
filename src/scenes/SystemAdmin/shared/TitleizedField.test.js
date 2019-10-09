import React from 'react';
import { shallow } from 'enzyme';

import TitleizedField from './TitleizedField';

describe('TitleizedField', () => {
  it('transforms text properly', () => {
    const record = { issuer: 'coast-guard' };
    const field = shallow(<TitleizedField source="issuer" record={record} />);

    expect(field.text()).toEqual('Coast Guard');
  });
});
