import React from 'react';
import { shallow } from 'enzyme';

import TitleizedField from './TitleizedField';

describe('TitleizedField', () => {
  it('transforms text properly', () => {
    const div = document.createElement('div');
    let record = { issuer: 'coast-guard' };
    let field = shallow(<TitleizedField source="issuer" record={record} />, div);

    expect(field.text()).toEqual('Coast Guard');

    record = { issuer: 'Army' };
    field = shallow(<TitleizedField source="issuer" record={record} />, div);

    expect(field.text()).toEqual('Army');
  });
});
