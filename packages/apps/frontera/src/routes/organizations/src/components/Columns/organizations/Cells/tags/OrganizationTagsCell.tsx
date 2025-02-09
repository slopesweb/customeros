import React, { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { Tag, DataSource } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions.ts';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';

import { TagsCell } from '../../../shared/Cells';

interface OrgCardProps {
  id: string;
}

export const OrganizationsTagsCell = observer(({ id }: OrgCardProps) => {
  const store = useStore();
  const [isEdit, setIsEdit] = useState(false);
  const organizationstore = store.organizations.value.get(id);
  const ref = useRef(null);

  useOutsideClick({
    ref: ref,
    handler: (e) => {
      // @ts-expect-error e.target.id can be undefined
      if (e?.target?.id.includes('react-select')) {
        e.preventDefault();
        e.stopPropagation();

        return;
      }
      setIsEdit(false);
    },
  });

  useEffect(() => {
    store.ui.setIsEditingTableCell(isEdit);
  }, [isEdit]);

  const handleCreateOption = (value: string) => {
    store.tags?.create({ name: value });

    organizationstore?.update((org) => {
      org.tags = [
        ...(org.tags || []),
        {
          id: value,
          name: value,
          metadata: {
            id: value,
            source: DataSource.Openline,
            sourceOfTruth: DataSource.Openline,
            appSource: 'organization',
            created: new Date().toISOString(),
            lastUpdated: new Date().toISOString(),
          },
          appSource: 'organization',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          source: DataSource.Openline,
        },
      ];

      return org;
    });
  };

  const handleChange = (tags: SelectOption<string>[]) => {
    organizationstore?.update((o) => {
      o.tags =
        (tags
          .map((tag) => store.tags?.value.get(tag.value)?.value)
          .filter(Boolean) as Array<Tag>) ?? [];

      return o;
    });
  };

  const tags = (organizationstore?.value?.tags ?? []).filter((d) => !!d?.name);

  return (
    <div ref={ref} onDoubleClick={() => setIsEdit(true)}>
      <TagsCell
        isEdit={isEdit}
        tags={tags ?? []}
        setIsEdit={setIsEdit}
        onChange={handleChange}
        onCreateOption={handleCreateOption}
      />
    </div>
  );
});
