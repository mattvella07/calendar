FROM scratch
COPY calendar-server docs.json docs.html /
COPY frontend/dist /frontend/dist
EXPOSE 8080
CMD ["./calendar-server"]