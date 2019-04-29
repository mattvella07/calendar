FROM scratch
ADD calendar-server /
ADD frontend/dist /frontend/dist
EXPOSE 8080
CMD ["./calendar-server"]