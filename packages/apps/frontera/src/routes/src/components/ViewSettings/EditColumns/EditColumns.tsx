import { useSearchParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import difference from 'lodash/difference';
import { observer } from 'mobx-react-lite';
import {
  Droppable,
  DragDropContext,
  OnDragEndResponder,
} from '@hello-pangea/dnd';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Columns02 } from '@ui/media/icons/Columns02';
import { Filter, TableViewType, ColumnViewType } from '@graphql/types';
import { Menu, MenuList, MenuGroup, MenuButton } from '@ui/overlay/Menu/Menu';
import { useTableColumnOptionsMap } from '@organizations/hooks/useTableColumnOptionsMap.tsx';

import { ColumnItem, DraggableColumnItem } from './ColumnItem';

interface EditColumnsProps {
  type: TableViewType;
}

export const EditColumns = observer(({ type }: EditColumnsProps) => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const preset = match(type)
    .with(
      TableViewType.Opportunities,
      () => store.tableViewDefs.opportunitiesPreset,
    )
    .otherwise(() => searchParams?.get('preset'));

  const [optionsMap, helperTextMap] = useTableColumnOptionsMap(type);

  const tableViewDef = store.tableViewDefs.getById(preset ?? '0');

  const columns =
    tableViewDef?.value?.columns.map((c) => ({
      ...c,
      label: optionsMap[c.columnType],
      helperText: helperTextMap[c.columnType],
    })) ?? [];

  const leadingPinnedColumns = match(tableViewDef?.value?.tableType)
    .with(TableViewType.Organizations, () =>
      columns.filter(({ columnType }) =>
        [
          ColumnViewType.OrganizationsAvatar,
          ColumnViewType.OrganizationsName,
        ].includes(columnType),
      ),
    )
    .with(TableViewType.Contacts, () =>
      columns.filter(({ columnType }) =>
        [ColumnViewType.ContactsAvatar, ColumnViewType.ContactsName].includes(
          columnType,
        ),
      ),
    )
    .with(TableViewType.Opportunities, () =>
      columns.filter((c) => {
        const filter = JSON.parse(c.filter) as Filter;
        const externalStage = filter?.AND?.find(
          (f) => f.filter?.property === 'externalStage',
        )?.filter?.value;

        return externalStage === 'STAGE1';
      }),
    )
    .otherwise(() => [columns[0]]);

  const traillingPinnedColumn = match(tableViewDef?.value?.tableType)
    .with(TableViewType.Opportunities, () =>
      columns.filter((c) => {
        const filter = JSON.parse(c.filter) as Filter;
        const internalStage = filter?.AND?.find(
          (f) => f.filter?.property === 'internalStage',
        )?.filter?.value;

        return ['CLOSED_WON', 'CLOSED_LOST'].includes(internalStage);
      }),
    )
    .otherwise(() => []);

  const draggableColumns = difference(
    columns,
    leadingPinnedColumns,
    traillingPinnedColumn,
  ).filter((d) => {
    return ![
      ColumnViewType.ContactsLanguages,
      ColumnViewType.ContactsSkills,
      ColumnViewType.ContactsSchools,
      ColumnViewType.ContactsExperience,
    ].includes(d?.columnType ?? '');
  });

  const handleDragEnd: OnDragEndResponder = (res) => {
    const sourceColumnId = draggableColumns[res.source.index]?.columnId;
    const destColumnId =
      draggableColumns[res?.destination?.index as number]?.columnId;
    const destination = res.destination;

    if (!destination) return;
    if (sourceColumnId === destColumnId) return;

    tableViewDef?.reorderColumn(sourceColumnId, destColumnId);
  };

  return (
    <>
      <Menu
        onOpenChange={(open) => {
          if (!open) {
            tableViewDef?.orderColumnsByVisibility();
          }
        }}
      >
        <MenuButton asChild>
          <Button size='xs' leftIcon={<Columns02 />} data-test={`edit-columns`}>
            Edit columns
          </Button>
        </MenuButton>
        <DragDropContext onDragEnd={handleDragEnd}>
          <MenuList className='w-[350px] max-h-[600px] overflow-y-auto'>
            {leadingPinnedColumns.map((col) => (
              <ColumnItem
                isPinned
                noPointerEvents
                visible={col?.visible}
                columnId={col?.columnId}
                columnType={col?.columnType}
                label={col?.name || col?.label}
                key={`${col?.columnType}-${col?.columnId}`}
              />
            ))}

            <Droppable
              key='active-columns'
              droppableId='active-columns'
              renderClone={(provided, snapshot, rubric) => {
                return (
                  <ColumnItem
                    provided={provided}
                    snapshot={snapshot}
                    visible={draggableColumns?.[rubric.source.index]?.visible}
                    columnId={draggableColumns?.[rubric.source.index]?.columnId}
                    helperText={
                      draggableColumns?.[rubric.source.index]?.helperText
                    }
                    columnType={
                      draggableColumns?.[rubric.source.index]?.columnType
                    }
                    label={
                      draggableColumns[rubric.source.index]?.name ||
                      draggableColumns[rubric.source.index]?.label
                    }
                    onCheck={(columnId) => {
                      tableViewDef?.update((value) => {
                        const columnIndex = value.columns.findIndex(
                          (c) => c.columnId === columnId,
                        );

                        value.columns[columnIndex].visible =
                          !value?.columns?.[columnIndex]?.visible;

                        return value;
                      });
                    }}
                  />
                );
              }}
            >
              {(provided, { isDraggingOver }) => (
                <>
                  <MenuGroup
                    ref={provided.innerRef}
                    {...provided.droppableProps}
                  >
                    {draggableColumns.map((col, index) => (
                      <DraggableColumnItem
                        index={index}
                        key={col?.columnType}
                        visible={col?.visible}
                        columnId={col?.columnId}
                        helperText={col?.helperText}
                        columnType={col?.columnType}
                        label={col?.name || col?.label}
                        noPointerEvents={isDraggingOver}
                        onCheck={(columnId) => {
                          tableViewDef?.update((value) => {
                            const columnIndex = value.columns.findIndex(
                              (c) => c.columnId === columnId,
                            );

                            value.columns[columnIndex].visible =
                              !value?.columns?.[columnIndex]?.visible;

                            return value;
                          });
                        }}
                      />
                    ))}
                    {provided.placeholder}
                  </MenuGroup>
                </>
              )}
            </Droppable>

            {traillingPinnedColumn.map((col) => (
              <ColumnItem
                isPinned
                noPointerEvents
                visible={col?.visible}
                columnId={col?.columnId}
                columnType={col?.columnType}
                label={col?.name || col?.label}
                key={`${col?.columnType}-${col?.columnId}`}
              />
            ))}
          </MenuList>
        </DragDropContext>
      </Menu>
    </>
  );
});
